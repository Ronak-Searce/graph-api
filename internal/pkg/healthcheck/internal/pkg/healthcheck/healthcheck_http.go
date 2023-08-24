package healthcheck

import (
	"encoding/json"
	"net/http"

	"github.com/Ronak-Searce/graph-api/internal/pkg/logger"
)

// ServeHTTP is an HTTP Handler.
func (hc *HealthCheck) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	w.Header().Set("Content-Type", "application/json")
	resp := hc.Run(ctx)

	if resp.Status != StatusError {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Errorf(ctx, "failed to encode failed healthcheck response: %v", err)
	}
}
