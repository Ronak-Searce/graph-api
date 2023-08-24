package tracing

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

var _ interface {
	graphql.HandlerExtension
	graphql.OperationInterceptor
} = graphqlInterceptor{}

type graphqlInterceptor struct{}

// GraphqlInterceptor returns graphql extension that adds some graphql metadata to span attributes.
func GraphqlInterceptor() graphql.HandlerExtension {
	return graphqlInterceptor{}
}

func (graphqlInterceptor) ExtensionName() string {
	return "tracing-interceptor"
}

func (graphqlInterceptor) Validate(_ graphql.ExecutableSchema) error {
	return nil
}

func (graphqlInterceptor) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	if graphql.HasOperationContext(ctx) {
		rc := graphql.GetOperationContext(ctx)
		if rc != nil {
			ctx = AddGraphqlToSpan(ctx, rc)
		}
	}

	return next(ctx)
}
