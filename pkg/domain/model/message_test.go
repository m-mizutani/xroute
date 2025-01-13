package model_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/m-mizutani/xroute/pkg/domain/model"
)

func TestNewGoogleIDToken(t *testing.T) {
	tests := []struct {
		name string
		src  map[string]any
		want *model.GoogleIDToken
	}{
		{
			name: "valid input",
			src: map[string]any{
				"aud":            []string{"test-aud"},
				"azp":            "test-azp",
				"email":          "test@example.com",
				"email_verified": true,
				"exp":            time.Unix(1234567890, 0),
				"iat":            time.Unix(1234567890, 0),
				"iss":            "test-iss",
				"sub":            "test-sub",
			},
			want: &model.GoogleIDToken{
				Aud:           []string{"test-aud"},
				Azp:           "test-azp",
				Email:         "test@example.com",
				EmailVerified: true,
				Exp:           time.Unix(1234567890, 0),
				Iat:           time.Unix(1234567890, 0),
				Iss:           "test-iss",
				Sub:           "test-sub",
			},
		},
		{
			name: "nil input",
			src:  nil,
			want: nil,
		},
		{
			name: "missing fields",
			src: map[string]any{
				"aud": []string{"test-aud"},
			},
			want: &model.GoogleIDToken{
				Aud: []string{"test-aud"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := model.NewGoogleIDToken(tt.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGoogleIDToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
