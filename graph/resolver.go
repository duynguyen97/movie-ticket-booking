package graph

import (
	"movie-ticket-booking/internal/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	authService  *services.AuthService
	movieService *services.MovieService
}

func NewResolver(authService *services.AuthService, movieService *services.MovieService) *Resolver {
	return &Resolver{
		authService:  authService,
		movieService: movieService,
	}
}
