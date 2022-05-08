package models

import "github.com/google/uuid"

// TestResponse implements the TestResponse GraphQL type
type TestResponse struct {
	uuid string
}

// NewTestResponse creates a new TestResponse instance
func NewTestResponse() *TestResponse {
	return &TestResponse{uuid: uuid.NewString()}
}

// UUID returns the uuid field
func (r *TestResponse) UUID() string {
	return r.uuid
}
