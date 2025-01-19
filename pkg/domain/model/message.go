package model

import "time"

// Message is general message schema for every message type.
type Message struct {
	// Source is origin of the message. E.g. "pubsub", "sns", "github", "raw", etc.
	Source string `json:"source"`

	// Schema is message schema. It's used to identify the message type.
	Schema string `json:"schema"`

	// Header is HTTP header of the message.
	Header map[string]string `json:"header"`

	// Body is parsed raw HTTP body. If content-type is application/json, it's parsed as JSON. Otherwise, it's raw bytes.
	Body any `json:"body"`

	// Data is parsed data part of the message. It's free format and can be any type. If it's JSON, it's parsed as JSON. Otherwise, it's raw bytes.
	Data any `json:"data"`

	// Auth is authentication information of the message. Authentication message is extracted from header mainly, e.g. JWT token, Secret key, etc.
	Auth AuthContext `json:"auth"`
}

// AuthContext is authentication context of the message.
type AuthContext struct {
	// Google is parsed Google ID Token. It's set if the message is authenticated by Google ID Token.
	Google *GoogleIDToken `json:"google,omitempty"`

	// GitHub is parsed GitHub authentication information.
	GitHub *AuthContextGitHub `json:"github,omitempty"`
}

// AuthContextGitHub is GitHub authentication information.
type AuthContextGitHub struct {
	Webhook *GitHubWebhookAuth `json:"webhook"`
	Actions map[string]any     `json:"actions"`
}

// GitHubActionsIDToken is parsed GitHub Actions ID Token.
type GitHubActionsIDToken struct {
	/* Example of GitHub Actions ID Token
	https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/about-security-hardening-with-openid-connect#understanding-the-oidc-token

	"jti": "example-id",
	"sub": "repo:octo-org/octo-repo:environment:prod",
	"environment": "prod",
	"aud": "https://github.com/octo-org",
	"ref": "refs/heads/main",
	"sha": "example-sha",
	"repository": "octo-org/octo-repo",
	"repository_owner": "octo-org",
	"actor_id": "12",
	"repository_visibility": "private",
	"repository_id": "74",
	"repository_owner_id": "65",
	"run_id": "example-run-id",
	"run_number": "10",
	"run_attempt": "2",
	"runner_environment": "github-hosted"
	"actor": "octocat",
	"workflow": "example-workflow",
	"head_ref": "",
	"base_ref": "",
	"event_name": "workflow_dispatch",
	"ref_type": "branch",
	"job_workflow_ref": "octo-org/octo-automation/.github/workflows/oidc.yml@refs/heads/main",
	"iss": "https://token.actions.githubusercontent.com",
	"nbf": 1632492967,
	"exp": 1632493867,
	"iat": 1632493567
	*/

	JTI string `json:"jti"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Ref string `json:"ref"`
	Sha string `json:"sha"`

	Repository           string `json:"repository"`
	RepositoryOwner      string `json:"repository_owner"`
	ActorID              string `json:"actor_id"`
	RepositoryVisibility string `json:"repository_visibility"`
	RepositoryID         string `json:"repository_id"`
	RepositoryOwnerID    string `json:"repository_owner_id"`
	RunID                string `json:"run_id"`
	RunNumber            string `json:"run_number"`
	RunAttempt           string `json:"run_attempt"`
	RunnerEnvironment    string `json:"runner_environment"`
	Actor                string `json:"actor"`
	Workflow             string `json:"workflow"`
	HeadRef              string `json:"head_ref"`
	BaseRef              string `json:"base_ref"`
	EventName            string `json:"event_name"`
	RefType              string `json:"ref_type"`
	JobWorkflowRef       string `json:"job_workflow_ref"`

	Iss string    `json:"iss"`
	Nbf time.Time `json:"nbf"`
	Exp time.Time `json:"exp"`
	Iat time.Time `json:"iat"`
}

func NewGitHubActionsIDToken(t map[string]any) *GitHubActionsIDToken {
	if t == nil {
		return nil
	}

	dst := &GitHubActionsIDToken{}
	for key, value := range t {
		switch key {
		case "jti":
			dst.JTI = value.(string)
		case "sub":
			dst.Sub = value.(string)
		case "aud":
			dst.Aud = value.(string)
		case "ref":
			dst.Ref = value.(string)
		case "sha":
			dst.Sha = value.(string)

		case "repository":
			dst.Repository = value.(string)
		case "repository_owner":
			dst.RepositoryOwner = value.(string)
		case "actor_id":
			dst.ActorID = value.(string)
		case "repository_visibility":
			dst.RepositoryVisibility = value.(string)
		case "repository_id":
			dst.RepositoryID = value.(string)
		case "repository_owner_id":
			dst.RepositoryOwnerID = value.(string)
		case "run_id":
			dst.RunID = value.(string)
		case "run_number":
			dst.RunNumber = value.(string)
		case "run_attempt":
			dst.RunAttempt = value.(string)
		case "runner_environment":
			dst.RunnerEnvironment = value.(string)
		case "actor":
			dst.Actor = value.(string)
		case "workflow":
			dst.Workflow = value.(string)
		case "head_ref":
			dst.HeadRef = value.(string)
		case "base_ref":
			dst.BaseRef = value.(string)
		case "event_name":
			dst.EventName = value.(string)
		case "ref_type":
			dst.RefType = value.(string)
		case "job_workflow_ref":
			dst.JobWorkflowRef = value.(string)

		case "iss":
			dst.Iss = value.(string)
		case "nbf":
			dst.Nbf = value.(time.Time)
		case "exp":
			dst.Exp = value.(time.Time)
		case "iat":
			dst.Iat = value.(time.Time)
		}
	}

	return dst
}

// GitHubWebhookAuth is parsed GitHub Webhook authentication information.
type GitHubWebhookAuth struct {
	// HookID is from "X-GitHub-Hook-ID" header.
	HookID int64 `json:"hook_id"`

	// TargetID is from "X-GitHub-Hook-Target-ID" header.
	TargetID int64 `json:"target_id"`

	// TargetType is from "X-GitHub-Hook-Target-Type" header.
	TargetType string `json:"target_type"`

	// Valid is true if the webhook is validated by secret key.
	Valid bool `json:"valid"`
}

// GoogleIDToken is parsed Google ID Token.
type GoogleIDToken struct {
	Aud           []string  `json:"aud"`
	Azp           string    `json:"azp"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Exp           time.Time `json:"exp"`
	Iat           time.Time `json:"iat"`
	Iss           string    `json:"iss"`
	Sub           string    `json:"sub"`
}

func NewGoogleIDToken(src map[string]any) *GoogleIDToken {
	if src == nil {
		return nil
	}

	dst := &GoogleIDToken{}
	for key, value := range src {
		switch key {
		case "aud":
			dst.Aud = value.([]string)
		case "azp":
			dst.Azp = value.(string)
		case "email":
			dst.Email = value.(string)
		case "email_verified":
			dst.EmailVerified = value.(bool)
		case "exp":
			dst.Exp = value.(time.Time)
		case "iat":
			dst.Iat = value.(time.Time)
		case "iss":
			dst.Iss = value.(string)
		case "sub":
			dst.Sub = value.(string)
		}
	}

	return dst
}
