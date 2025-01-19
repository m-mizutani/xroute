package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/xroute/pkg/domain/interfaces"
	"github.com/m-mizutani/xroute/pkg/domain/model"
)

func validateGitHubActionToken(ctx context.Context, authHdr string) (map[string]any, error) {
	hdr := strings.SplitN(authHdr, " ", 2)

	// Skip if not Bearer token
	if len(hdr) != 2 || strings.ToLower(hdr[0]) != "bearer" {
		return nil, nil
	}

	jwksURL := "https://token.actions.githubusercontent.com/.well-known/jwks"
	token := hdr[1]

	set, err := jwk.Fetch(ctx, jwksURL)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to fetch JWK set")
	}

	parsed, err := jwt.ParseString(token, jwt.WithKeySet(set))
	if err != nil {
		return nil, goerr.Wrap(err, "failed to parse JWT token as GitHub Action token", goerr.V("token", trimToken(hdr[1])))
	}

	claims, err := parsed.AsMap(ctx)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to convert JWT token to map", goerr.V("token", trimToken(token)))
	}

	return claims, nil
}

func handleGitHubActions(r *http.Request, uc interfaces.UseCases) error {
	ctx := r.Context()
	authHdr := r.Header.Get("Authorization")
	claims, err := validateGitHubActionToken(ctx, authHdr)
	if err != nil {
		return err
	}
	if claims == nil {
		return nil
	}

	var payload any
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return goerr.Wrap(err, "failed to read HTTP body")
	}

	switch r.Header.Get("Content-Type") {
	case "application/json":
		if err := json.Unmarshal(body, &payload); err != nil {
			return goerr.Wrap(err, "failed to parse JSON body", goerr.V("body", string(body)))
		}
	default:
		payload = body
	}

	msg := model.Message{
		Source: "github.actions",
		Schema: "actions",
		Header: cloneHeader(r.Header),
		Body:   body,
		Data:   payload,
		Auth: model.AuthContext{
			GitHub: &model.AuthContextGitHub{
				Actions: claims,
			},
		},
	}

	if err := uc.Route(ctx, msg); err != nil {
		return goerr.Wrap(err, "failed to route message", goerr.V("msg", msg))
	}

	return nil
}
