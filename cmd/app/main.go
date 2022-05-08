package main

import (
	"log"

	"github.com/aidenwallis/api-gateway-cdk-graphql-test/internal/handler"
	"github.com/aidenwallis/api-gateway-cdk-graphql-test/internal/resolvers"
	"github.com/aidenwallis/api-gateway-cdk-graphql-test/schema"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/graph-gophers/graphql-go"
)

func main() {
	parsedSchema, err := schema.String()
	if err != nil {
		log.Fatalln("Failed to init schema: " + err.Error())
	}

	rootResolver := resolvers.NewRootResolver()

	h := &handler.Handler{
		Schema: graphql.MustParseSchema(parsedSchema, rootResolver),
	}

	lambda.Start(h.HandleRequest)
}
