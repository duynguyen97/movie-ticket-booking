package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.66

import (
	"context"
	"fmt"
	"movie-ticket-booking/graph/generated"
)

// Dummy is the resolver for the dummy field.
func (r *mutationResolver) Dummy(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented: Dummy - dummy"))
}

// Ping is the resolver for the ping field.
func (r *queryResolver) Ping(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented: Ping - ping"))
}

// Dummy is the resolver for the dummy field.
func (r *subscriptionResolver) Dummy(ctx context.Context) (<-chan string, error) {
	panic(fmt.Errorf("not implemented: Dummy - dummy"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
