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
	"github.com/habinkwon/chat-app/graph"
	"github.com/habinkwon/chat-app/graph/generated"
	"github.com/habinkwon/chat-app/pkg/middleware/auth"
	mysqlrepo "github.com/habinkwon/chat-app/pkg/repository/mysql"
	redisrepo "github.com/habinkwon/chat-app/pkg/repository/redis"
	"github.com/habinkwon/chat-app/pkg/service"
	"github.com/rs/cors"
)

func main() {
	listenAddr := flag.String("listen", ":8080", "")
	mysqlAddr := flag.String("mysql", "root@tcp(maria)/chat?parseTime=true", "")
	redisAddr := flag.String("redis", "redis:6379", "")
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
		UserSvc: &service.User{},
		ChatSvc: &service.Chat{
			IDNode:          idNode,
			ChatRepo:        &mysqlrepo.Chat{DB: db},
			ChatMemberRepo:  &mysqlrepo.ChatMember{DB: db},
			ChatMessageRepo: &mysqlrepo.ChatMessage{DB: db},
			ChannelRepo:     &redisrepo.Channel{Redis: rdb},
		},
	}
	hs := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler)
	r.Use(auth.Middleware())
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", hs)

	l, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatal(fmt.Errorf("error listening %s: %w", *listenAddr, err))
	}
	log.Printf("server listening on %s", *listenAddr)

	srv := &http.Server{
		Handler: r,
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
