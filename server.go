package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/neomarica/undergraduate-project/graph"
	"github.com/neomarica/undergraduate-project/graph/generated"
)

func main() {
	listen := flag.String("listen", ":8080", "")
	flag.Parse()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("server listening on %s", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
