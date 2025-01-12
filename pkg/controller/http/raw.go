package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/transmith/pkg/domain/interfaces"
	"github.com/m-mizutani/transmith/pkg/domain/model"
	"github.com/m-mizutani/transmith/pkg/utils/logging"
)

func handleRawMessage(r *http.Request, uc interfaces.UseCases) error {
	ctx := r.Context()
	logger := logging.Extract(ctx)
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return goerr.Wrap(err, "Unable to read request body")
	}

	// Build message from generic HTTP request
	msg := model.Message{
		Header: map[string]string{},
		Schema: r.URL.Query().Get("schema"),
	}

	// Copy HTTP headers to message header. Only the first value is stored.
	for k, v := range r.Header {
		msg.Header[k] = v[0]
	}

	// Unmarshal request body as JSON if Content-Type is application/json
	if r.Header.Get("Content-Type") == "application/json" {
		var data any
		if err := json.Unmarshal(raw, &data); err != nil {
			return goerr.Wrap(err, "Failed to unmarshal JSON", goerr.V("data", string(raw)))
		}
		msg.Data = data
		logger.Debug("Parsed data of Pub/Sub as JSON", "data", data)
	} else {
		msg.Data = string(raw)
		logger.Debug("Parsed data of Pub/Sub as string", "data", string(raw))
	}

	if err := uc.Transmit(ctx, msg); err != nil {
		return err
	}

	return nil
}
