package mutations

import "github.com/aidenwallis/api-gateway-cdk-graphql-test/internal/resolvers/models"

// Mutations implements mutation resolvers
type Mutations struct{}

// TestMutation implements the testMutation resolver
func (r *Mutations) TestMutation() *models.TestResponse {
	return models.NewTestResponse()
}
