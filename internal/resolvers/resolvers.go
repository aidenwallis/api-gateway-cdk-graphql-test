package resolvers

import (
	"github.com/aidenwallis/api-gateway-cdk-graphql-test/internal/resolvers/mutations"
	"github.com/aidenwallis/api-gateway-cdk-graphql-test/internal/resolvers/query"
)

// RootResolver implements the root resolver
type RootResolver struct {
	*query.Query
	*mutations.Mutations
}

// NewRootResolver creates a new instance of RootResolver
func NewRootResolver() *RootResolver {
	return &RootResolver{
		Query:     &query.Query{},
		Mutations: &mutations.Mutations{},
	}
}
