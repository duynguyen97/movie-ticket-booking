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
	"movie-ticket-booking/internal/middleware"
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

	// Initialize Redis client
	redisClient, err := database.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize services
	authService := services.NewAuthService(postgresDB.DB, "your-jwt-secret-key") // Replace with actual secret from config
	movieService := services.NewMovieService(postgresDB.DB)
	bookingService := services.NewBookingService(postgresDB.DB, redisClient.Client)

	// Create resolver with services
	resolver := graph.NewResolver(authService, movieService, bookingService)

	// Create GraphQL server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", middleware.AuthMiddleware(authService)(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}