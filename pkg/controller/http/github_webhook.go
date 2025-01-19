package http

import (
	"net/http"
	"strconv"

	"github.com/google/go-github/v68/github"
	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/xroute/pkg/domain/interfaces"
	"github.com/m-mizutani/xroute/pkg/domain/model"
)

func handleGitHubWebhook(r *http.Request, uc interfaces.UseCases, secret string) error {
	var secretValidated bool
	payload, err := github.ValidatePayload(r, []byte(secret))
	if err != nil {
		return goerr.Wrap(err, "Failed to validate GitHub webhook payload")
	}
	if secret != "" {
		secretValidated = true
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return goerr.Wrap(err, "Failed to parse GitHub webhook event", goerr.V("payload", string(payload)))
	}

	msg := model.Message{
		Source: "github.webhook",
		Schema: r.Header.Get("X-GitHub-Event"),
		Data:   event,
		Body:   payload,
		Auth: model.AuthContext{
			GitHub: &model.AuthContextGitHub{
				Webhook: &model.GitHubWebhookAuth{
					Valid: secretValidated,
				},
			},
		},
	}

	if v := r.Header.Get("X-GitHub-Hook-ID"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return goerr.Wrap(err, "Failed to parse X-GitHub-Hook-ID", goerr.V("value", v))
		}
		msg.Auth.GitHub.Webhook.HookID = id
	}

	if v := r.Header.Get("X-GitHub-Hook-Target-ID"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return goerr.Wrap(err, "Failed to parse X-GitHub-Hook-Target-ID", goerr.V("value", v))
		}
		msg.Auth.GitHub.Webhook.TargetID = id
	}

	if v := r.Header.Get("X-GitHub-Hook-Target-Type"); v != "" {
		msg.Auth.GitHub.Webhook.TargetType = v
	}

	if err := uc.Route(r.Context(), msg); err != nil {
		return goerr.Wrap(err, "Failed to route message")
	}

	return nil
}
