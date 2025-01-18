package http_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"testing"

	"github.com/google/go-github/v68/github"
	"github.com/m-mizutani/gt"
	http_server "github.com/m-mizutani/xroute/pkg/controller/http"
	"github.com/m-mizutani/xroute/pkg/domain/model"
	"github.com/m-mizutani/xroute/pkg/mock"
)

func TestHandleGitHubWebhook(t *testing.T) {
	secret := "my_secret"
	payload := []byte(`{"action":"opened"}`)
	signature := "sha256=" + generateSignature(payload, secret)

	tests := []struct {
		name           string
		headers        map[string]string
		secret         string
		payload        []byte
		expectedError  bool
		expectedSecret bool
	}{
		{
			name: "valid payload with secret",
			headers: map[string]string{
				"Content-Type":              "application/json",
				"X-GitHub-Event":            "issues",
				"X-GitHub-Hook-ID":          "12345",
				"X-GitHub-Hook-Target-ID":   "67890",
				"X-GitHub-Hook-Target-Type": "repository",
				"X-Hub-Signature-256":       signature,
			},
			secret:         secret,
			payload:        payload,
			expectedError:  false,
			expectedSecret: true,
		},
		{
			name: "invalid payload signature",
			headers: map[string]string{
				"Content-Type":              "application/json",
				"X-GitHub-Event":            "issues",
				"X-GitHub-Hook-ID":          "12345",
				"X-GitHub-Hook-Target-ID":   "67890",
				"X-GitHub-Hook-Target-Type": "repository",
				"X-Hub-Signature-256":       "invalidsignature",
			},
			secret:         secret,
			payload:        payload,
			expectedError:  true,
			expectedSecret: false,
		},
		{
			name: "no secret provided",
			headers: map[string]string{
				"Content-Type":              "application/json",
				"X-GitHub-Event":            "issues",
				"X-GitHub-Hook-ID":          "12345",
				"X-GitHub-Hook-Target-ID":   "67890",
				"X-GitHub-Hook-Target-Type": "repository",
			},
			secret:         "",
			payload:        payload,
			expectedError:  false,
			expectedSecret: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/", bytes.NewBuffer(tt.payload))
			gt.NoError(t, err)

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			mockUC := &mock.UseCasesMock{
				RouteFunc: func(ctx context.Context, msg model.Message) error {
					return nil
				},
			}

			err = http_server.HandleGitHubWebhook(req, mockUC, tt.secret)
			if tt.expectedError {
				gt.Error(t, err)
			} else {
				gt.NoError(t, err)

				gt.A(t, mockUC.RouteCalls()).Length(1).At(0, func(t testing.TB, v struct {
					Ctx context.Context
					Msg model.Message
				}) {
					if tt.secret != "" {
						gt.NotEqual(t, v.Msg.Auth.GitHub, nil)
						gt.NotEqual(t, v.Msg.Auth.GitHub.Webhook, nil)

						auth := v.Msg.Auth.GitHub.Webhook
						gt.Equal(t, auth.HookID, 12345)
						gt.Equal(t, auth.TargetID, 67890)
						gt.Equal(t, auth.TargetType, "repository")
						gt.Equal(t, auth.Valid, true)
					}

					data := v.Msg.Data.(*github.IssuesEvent)
					gt.NotEqual(t, data, nil)
					gt.NotEqual(t, data.Action, nil)
					gt.Equal(t, *data.Action, "opened")
				})
			}
		})
	}
}

func generateSignature(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	signature := hex.EncodeToString(h.Sum(nil))

	return signature
}
