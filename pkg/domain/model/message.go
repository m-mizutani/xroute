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
