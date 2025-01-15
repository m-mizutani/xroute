package safe_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/go-test/deep"
	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/xroute/pkg/utils/logging"
	"github.com/m-mizutani/xroute/pkg/utils/safe"
)

type mockWriter struct {
	writeErr error
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	return len(p), nil
}

type mockCloser struct {
	closeErr error
}

func (m *mockCloser) Close() error {
	return m.closeErr
}

func TestWrite(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer

	// Test successful write
	safe.Write(ctx, &buf, []byte("test data"))
	if diff := deep.Equal("test data", buf.String()); diff != nil {
		t.Error(diff)
	}

	// Test write error
	mockW := &mockWriter{writeErr: errors.New("write error")}
	log := slog.New(slog.NewJSONHandler(&buf, nil))
	ctx = logging.Inject(ctx, log)
	safe.Write(ctx, mockW, []byte("test data"))
	gt.S(t, buf.String()).Contains("write error")
}

func TestClose(t *testing.T) {
	ctx := context.Background()

	// Test successful close
	mockC := &mockCloser{}
	safe.Close(ctx, mockC)

	// Test close error
	mockC = &mockCloser{closeErr: errors.New("close error")}
	var buf bytes.Buffer
	log := slog.New(slog.NewJSONHandler(&buf, nil))
	ctx = logging.Inject(ctx, log)
	safe.Close(ctx, mockC)
	gt.S(t, buf.String()).Contains("close error")
}
