package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"movie-ticket-booking/graph"
	"movie-ticket-booking/graph/generated"
	"movie-ticket-booking/internal/config"
	"movie-ticket-booking/internal/database"
	"movie-ticket-booking/internal/services"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize database connection
	postgresDB, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize auth service
	authService := services.NewAuthService(postgresDB.DB, "your-jwt-secret-key") // Replace with actual secret from config

	// Create resolver with auth service
	resolver := graph.NewResolver(authService)

	// Create GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}