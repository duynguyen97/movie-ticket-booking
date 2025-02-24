package graph

import (
	"movie-ticket-booking/internal/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	authService    *services.AuthService
	movieService  *services.MovieService
	bookingService *services.BookingService
}

func NewResolver(authService *services.AuthService, movieService *services.MovieService, bookingService *services.BookingService) *Resolver {
	return &Resolver{
		authService:    authService,
		movieService:   movieService,
		bookingService: bookingService,
	}
}
