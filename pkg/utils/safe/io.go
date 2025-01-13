package safe

import (
	"context"
	"io"

	"github.com/m-mizutani/xroute/pkg/utils/logging"
)

func Write(ctx context.Context, w io.Writer, data []byte) {
	if _, err := w.Write(data); err != nil {
		logging.Extract(ctx).Error("Failed to write data", "error", err)
	}
}

func Close(ctx context.Context, closer io.Closer) {
	if err := closer.Close(); err != nil {
		logging.Extract(ctx).Error("Failed to close", "error", err)
	}
}
