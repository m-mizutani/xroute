package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/xroute/pkg/domain/interfaces"
	"github.com/m-mizutani/xroute/pkg/domain/model"
	"github.com/m-mizutani/xroute/pkg/domain/types"
	"github.com/m-mizutani/xroute/pkg/utils/logging"
)

type pubsubMessage struct {
	Message struct {
		Attributes  map[string]string `json:"attributes"`
		Data        []byte            `json:"data"`
		MessageID   string            `json:"message_id"`
		PublishTime time.Time         `json:"publish_time"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func handlePubSubMessage(r *http.Request, uc interfaces.UseCases) error {
	ctx := r.Context()
	logger := logging.Extract(ctx)
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return goerr.Wrap(err, "Unable to read request body")
	}

	// Extract Data part from Pub/Sub message
	var pubsubMsg pubsubMessage
	err = json.Unmarshal(raw, &pubsubMsg)
	if err != nil {
		return goerr.Wrap(err, "Invalid JSON",
			goerr.V("body", string(raw)),
			goerr.T(types.ErrTagBadRequest),
		)
	}
	logger.Debug("Received Pub/Sub event", "msg", pubsubMsg)

	// Build message from HTTP request and Pub/Sub message
	msg := model.Message{
		Header: map[string]string{},
		Schema: r.PathValue("schema"),
	}

	// Copy HTTP headers to message header. Only the first value is stored.
	for k, v := range r.Header {
		msg.Header[k] = v[0]
	}

	// Pub/Sub message body must be parsable as JSON
	var body any
	if err := json.Unmarshal(raw, &body); err != nil {
		return goerr.Wrap(err, "Failed to unmarshal JSON",
			goerr.V("data", string(raw)),
			goerr.T(types.ErrTagBadRequest),
		)
	}
	msg.Body = body

	// Pub/Sub message data is free format. If it can be parsed as JSON, it will be stored in msg.Data
	var data any
	if err := json.Unmarshal(pubsubMsg.Message.Data, &data); err == nil {
		msg.Data = data
		logger.Debug("Parsed data of Pub/Sub as JSON", "data", data)
	} else {
		msg.Data = string(pubsubMsg.Message.Data)
		logger.Debug("Data of Pub/Sub can not be parsed, use it as raw", "data", msg.Data)
	}

	// Extract Google ID token from Authorization header. Empty Authorization header is allowed, but ID token validation error is not allowed and return error.
	if authHdr := r.Header.Get("Authorization"); authHdr != "" {
		token, err := validateGoogleIDToken(ctx, authHdr)
		if err != nil {
			return goerr.Wrap(err, "Failed to validate Google ID token", goerr.T(types.ErrTagUnauthorized))
		}
		if token != nil {
			msg.Auth.Google = token
		}
	}

	if err := uc.Route(ctx, msg); err != nil {
		return err
	}

	return nil
}

func validateGoogleIDToken(ctx context.Context, authHdr string) (*model.GoogleIDToken, error) {
	hdr := strings.SplitN(authHdr, " ", 2)

	// Skip if not Bearer token
	if len(hdr) != 2 || strings.ToLower(hdr[0]) != "bearer" {
		return nil, nil
	}

	jwksURL := "https://www.googleapis.com/oauth2/v3/certs"

	set, err := jwk.Fetch(ctx, jwksURL)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to fetch JWK set")
	}

	token, err := jwt.ParseString(hdr[1], jwt.WithKeySet(set))
	if err != nil {
		return nil, goerr.Wrap(err, "failed to parse JWT token as Google ID Token", goerr.V("token", trimToken(hdr[1])))
	}

	claims, err := token.AsMap(ctx)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to convert JWT token to map", goerr.V("token", trimToken(hdr[1])))
	}

	return model.NewGoogleIDToken(claims), nil
}
