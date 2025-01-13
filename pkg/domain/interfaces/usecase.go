package interfaces

import (
	"context"

	"github.com/m-mizutani/xroute/pkg/domain/model"
)

type UseCases interface {
	Transmit(ctx context.Context, msg model.Message) error
}
