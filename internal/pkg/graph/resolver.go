package graph

import (
	"graph-api/pkg/graph"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	res IResolver
}

// IResolver ...
type IResolver interface {
	graph.QueryResolver
	graph.MutationResolver
	// graph.SubscriptionResolver
}

// InitImpl ...
func (r *Resolver) InitImpl(res IResolver) {
	r.res = res
}
