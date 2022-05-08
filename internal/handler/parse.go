package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

// Source: https://github.com/tonyghita/graphql-go-example/blob/main/handler/graphql.go
// A request respresents an HTTP request to the GraphQL endpoint.
// A request can have a single query or a batch of requests with one or more queries.
// It is important to distinguish between a single query request and a batch request with a single query.
// The shape of the response will differ in both cases.
type request struct {
	queries []query
	isBatch bool
}

// A query represents a single GraphQL query.
type query struct {
	OpName    string                 `json:"operationName"`
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

var errMethod = errors.New("unsupported method")

// This code is largely sourced from https://github.com/tonyghita/graphql-go-example/blob/main/handler/parse.go
// - but has been adapted slightlyu to work with lamdbas weird HTTP proxy types

func parse(req *events.APIGatewayProxyRequest) (*request, error) {
	switch req.HTTPMethod {
	case http.MethodPost:
		return parsePost(req.Body), nil
	case http.MethodGet:
		return parseGet(req.MultiValueQueryStringParameters), nil
	}

	return nil, errMethod
}

func parseGet(v url.Values) *request {
	var (
		queries   = v["query"]
		names     = v["operationName"]
		variables = v["variables"]
		qLen      = len(queries)
		nLen      = len(names)
		vLen      = len(variables)
	)

	if qLen == 0 {
		return &request{}
	}

	var requests = make([]query, 0, qLen)
	var isBatch bool

	// This loop assumes there will be a corresponding element at each index
	// for query, operation name, and variable fields.
	//
	// NOTE: This could be a bad assumption. Maybe we want to do some validation?
	for i, q := range queries {
		var n string
		if i < nLen {
			n = names[i]
		}

		var m = map[string]interface{}{}
		if i < vLen {
			str := variables[i]
			if err := json.Unmarshal([]byte(str), &m); err != nil {
				m = nil // TODO: Improve error handling here.
			}
		}

		requests = append(requests, query{Query: q, OpName: n, Variables: m})
	}

	if qLen > 1 {
		isBatch = true
	}

	return &request{queries: requests, isBatch: isBatch}
}

func parsePost(body string) *request {
	if len(body) == 0 {
		return &request{}
	}

	var queries []query
	var isBatch bool

	// Inspect the first character to inform how the body is parsed.
	switch body[0] {
	case '{':
		q := query{}
		err := json.Unmarshal([]byte(body), &q)
		if err == nil {
			queries = append(queries, q)
		}
	case '[':
		isBatch = true
		_ = json.Unmarshal([]byte(body), &queries)
	}

	return &request{queries: queries, isBatch: isBatch}
}
