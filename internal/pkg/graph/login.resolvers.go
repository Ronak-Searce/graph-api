package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"
	"graph-api/pkg/graph"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input graph.LoginInput) (*graph.LoginOutput, error) {
	return r.res.Login(ctx, input) //r.implementation.Login(ctx,input)
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
