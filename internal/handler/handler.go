package handler

import (
	"context"
	"net/http"
	"sync"

	"github.com/aidenwallis/api-gateway-cdk-graphql-test/internal/errors"
	"github.com/aidenwallis/go-utils/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/graph-gophers/graphql-go"
)

type Handler struct {
	Schema *graphql.Schema
}

func (h *Handler) HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req, err := parse(&event)
	if err != nil {
		return res(utils.Ternary(err == errMethod, http.StatusMethodNotAllowed, http.StatusInternalServerError)).body(err.Error()).end(), nil
	}

	if len(req.queries) == 0 {
		// nothing to do
		return res(http.StatusBadRequest).end(), nil
	}

	var (
		responses = make([]*graphql.Response, len(req.queries))
		wg        sync.WaitGroup
	)

	wg.Add(len(req.queries))

	for i, q := range req.queries {
		// Iterate through the parsed queries from the request.
		// These queries are executed in separate goroutines so they process in parallel.
		go func(i int, q query) {
			defer wg.Done()

			res := h.Schema.Exec(ctx, q.Query, q.OpName, q.Variables)

			// We have to do some work here to expand errors when it is possible for a resolver to return
			// more than one error (for example, a list resolver).
			res.Errors = errors.Expand(res.Errors)

			responses[i] = res
		}(i, q)
	}

	wg.Wait()

	// json takes interface{} anyways, so using it in this generic is ifne
	return res(http.StatusOK).json(utils.Ternary[interface{}](req.isBatch, responses, responses[0])).end(), nil
}
