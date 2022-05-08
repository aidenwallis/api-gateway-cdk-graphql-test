package query

import "github.com/aidenwallis/api-gateway-cdk-graphql-test/internal/resolvers/models"

// Query implements the query resolvers
type Query struct{}

// TestQuery implements testQuery
func (r *Query) TestQuery() *models.TestResponse {
	return models.NewTestResponse()
}
