package core

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
)

type application struct {
	graphPort uint32
	graph     *handler.Server

	srv *http.Server
}

// New creates core application
func New(_ context.Context, httpPort uint32) *application {
	a := &application{
		graphPort: httpPort,
	}

	return a
}
