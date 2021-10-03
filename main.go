package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bwmarrin/snowflake"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/neomarica/undergraduate-project/graph"
	"github.com/neomarica/undergraduate-project/graph/generated"
	"github.com/neomarica/undergraduate-project/pkg/middleware/auth"
	"github.com/neomarica/undergraduate-project/pkg/repository/mysql"
	"github.com/neomarica/undergraduate-project/pkg/service"
)

func main() {
	listenAddr := flag.String("listen", ":8080", "")
	mysqlAddr := flag.String("mysql", "root@/chat", "")
	redisAddr := flag.String("redis", "localhost:6379", "")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		sig := <-sigs
		log.Printf("caught %s", sig)
		cancel()
	}()

	db, err := sql.Open("mysql", *mysqlAddr)
	if err != nil {
		log.Fatal(fmt.Errorf("error opening mysql: %w", err))
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: *redisAddr,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal(fmt.Errorf("error opening redis: %w", err))
	}
	defer rdb.Close()

	idNode, err := snowflake.NewNode(0)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing snowflake: %w", err))
	}

	resolver := &graph.Resolver{
		Chat: &service.Chat{
			Chat:        &mysql.Chat{DB: db},
			ChatMember:  &mysql.ChatMember{DB: db},
			ChatMessage: &mysql.ChatMessage{DB: db},
			IDNode:      idNode,
		},
	}
	hs := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	mux := chi.NewMux()
	mux.Use(auth.Middleware())
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", hs)

	l, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatal(fmt.Errorf("error listening %s: %w", *listenAddr, err))
	}
	log.Printf("server listening on %s", *listenAddr)

	srv := &http.Server{
		Handler: mux,
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := srv.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Print(fmt.Errorf("server error: %w", err))
		}
	}()
	go func() {
		<-ctx.Done()
		if err := srv.Close(); err != nil {
			log.Print(fmt.Errorf("error closing server: %w", err))
		}
	}()
	<-done
}
