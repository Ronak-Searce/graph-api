package gqlclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/rs/zerolog/log"
)

// GQLClient ...
type GQLClient struct {
	cli         *http.Client
	accessToken string
	envID       string
}

// NewGQLClient ...
func NewGQLClient(cli *http.Client) *GQLClient {
	return &GQLClient{
		cli:   cli,
		envID: os.Getenv("ENV_ID"),
	}
}

// WithAccessToken ...
func (c *GQLClient) WithAccessToken(t string) *GQLClient {
	return &GQLClient{
		cli:         c.cli,
		accessToken: t,
		envID:       c.envID,
	}
}

// Do ...
func (c *GQLClient) Do(ctx context.Context, body string, responseTarget interface{}) error {
	return c.DoWithURL(ctx, "http://localhost:8080/query", body, responseTarget)
}

// DoWithURL ...
func (c *GQLClient) DoWithURL(ctx context.Context, url, body string, responseTarget interface{}) error {
	reqBody, err := json.Marshal(struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     body,
		Variables: nil,
	})
	if err != nil {
		return fmt.Errorf("marshal request: %v", err)
	}
	log.Debug().Str("gql", body).Msg("send request")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if t := c.accessToken; t != "" {
		req.Header.Set("Authorization", "Bearer "+t)
	}

	if envID := c.envID; envID != "" {
		req.Header.Set("env-id", envID)
	}

	response, err := c.cli.Do(req)
	if err != nil {
		return fmt.Errorf("do: %v", err)
	}
	defer response.Body.Close()
	log.Debug().Str("status", response.Status).Msg("got response")

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read response: %v", err)
	}
	log.Debug().RawJSON("body", responseBody).Msg("got response")

	var errorResponse graphql.Response
	err = json.Unmarshal(responseBody, &errorResponse)
	if err != nil {
		return fmt.Errorf("unmarshal error response: %v", err)
	}

	if len(errorResponse.Errors) > 0 {
		return fmt.Errorf("error response: %s", string(responseBody))
	}

	if err := json.Unmarshal(responseBody, &responseTarget); err != nil {
		return fmt.Errorf("unmarshal error response: %v", err)
	}

	return nil
}
