package handler

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// response is a simple wrapper to make this slightly less awful to work with, by adding chainable funcs
type response struct {
	events.APIGatewayProxyResponse
}

// res creates a new response
func res(statusCode int) *response {
	return &response{
		APIGatewayProxyResponse: events.APIGatewayProxyResponse{
			Headers: map[string]string{
				// a drawback of this setup is that we have to send these headers ourselves, because lambda proxies
				// don't get their headers intercepted.
				"access-control-allow-origin":  "*",
				"access-control-allow-headers": "*",
				"access-control-allow-methods": "GET, POST, OPTIONS",
			},
			StatusCode: statusCode,
		},
	}
}

// status sets the status
func (r *response) status(statusCode int) *response {
	r.APIGatewayProxyResponse.StatusCode = statusCode
	return r
}

// body writes a raw text body
func (r *response) body(body string) *response {
	r.APIGatewayProxyResponse.Body = body
	return r
}

// json is a shorthand to add a json body
func (r *response) json(v interface{}) *response {
	bs, _ := json.Marshal(v) // if you read this, please don't actually ignore errors like this
	r.APIGatewayProxyResponse.Body = string(bs)
	return r
}

// end returns the finished object
func (r *response) end() events.APIGatewayProxyResponse {
	return r.APIGatewayProxyResponse
}
