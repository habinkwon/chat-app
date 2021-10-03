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
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/neomarica/undergraduate-project/graph"
	"github.com/neomarica/undergraduate-project/graph/generated"
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

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", h)

	l, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatal(fmt.Errorf("error listening %s: %w", *listenAddr, err))
	}
	log.Printf("server listening on %s", *listenAddr)

	done := make(chan struct{})
	srv := &http.Server{
		Handler: h,
	}
	go func() {
		defer close(done)
		if err := srv.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Print(fmt.Errorf("error serving: %w", err))
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
