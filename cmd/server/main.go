package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"
	"github.com/nico-ulbricht/bullauge/pkg/pod"
)

func init() {
	godotenv.Load()
}

func main() {
	fields := graphql.Fields{
		"pods": &pod.Query,
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to generate Schema: %v", err)
	}

	gqlHandler := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	http.HandleFunc("/healthz", healthHandler)
	http.Handle("/graphql", gqlHandler)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("Listening on port %s.", port)
	http.ListenAndServe(port, nil)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}
