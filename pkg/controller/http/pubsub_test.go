package http_test

import (
	"bytes"
	"context"
	_ "embed"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/xroute/pkg/controller/http"
	"github.com/m-mizutani/xroute/pkg/domain/model"
	"github.com/m-mizutani/xroute/pkg/mock"
)

//go:embed testdata/pubsub_json.json
var pubsubJSON []byte

func TestPubSubJSON(t *testing.T) {
	uc := &mock.UseCasesMock{
		TransmitFunc: func(ctx context.Context, msg model.Message) error {
			return nil
		},
	}
	srv := http.New(uc)

	r := httptest.NewRequest("POST", "/msg/pubsub/json_schema", bytes.NewReader(pubsubJSON))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)

	gt.Equal(t, w.Code, 200)
	gt.A(t, uc.TransmitCalls()).Length(1).At(0, func(t testing.TB, v struct {
		Ctx context.Context
		Msg model.Message
	}) {
		msg := gt.Cast[map[string]any](t, v.Msg.Data)
		gt.Equal(t, msg["kind"], "storage#object")
		gt.Equal(t, v.Msg.Schema, "json_schema")
	})
}

//go:embed testdata/pubsub_text.json
var pubsubText []byte

func TestPubSubText(t *testing.T) {
	uc := &mock.UseCasesMock{
		TransmitFunc: func(ctx context.Context, msg model.Message) error {
			return nil
		},
	}
	srv := http.New(uc)

	r := httptest.NewRequest("POST", "/msg/pubsub/text_schema", bytes.NewReader(pubsubText))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)

	gt.Equal(t, w.Code, 200)
	gt.A(t, uc.TransmitCalls()).Length(1).At(0, func(t testing.TB, v struct {
		Ctx context.Context
		Msg model.Message
	}) {
		msg := gt.Cast[string](t, v.Msg.Data)
		gt.Equal(t, msg, "Hello, World")
		gt.Equal(t, v.Msg.Schema, "text_schema")
	})
}

func TestPubSubAuth(t *testing.T) {
	idToken, ok := os.LookupEnv("TEST_GOOGLE_ID_TOKEN")
	if !ok {
		t.Skip("TEST_GOOGLE_ID_TOKEN is not set")
	}
	email, ok := os.LookupEnv("TEST_GOOGLE_EMAIL")
	if !ok {
		t.Skip("TEST_GOOGLE_EMAIL is not set")
	}

	uc := &mock.UseCasesMock{}
	srv := http.New(uc)

	r := httptest.NewRequest("POST", "/msg/pubsub/json_schema", bytes.NewReader(pubsubJSON))
	r.Header.Set("Authorization", "Bearer "+idToken)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)

	gt.Equal(t, w.Code, 200)
	gt.A(t, uc.TransmitCalls()).Length(1).At(0, func(t testing.TB, v struct {
		Ctx context.Context
		Msg model.Message
	}) {
		gt.Equal(t, v.Msg.Auth.Google.Email, email)
	})
}
