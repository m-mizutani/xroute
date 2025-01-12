package interfaces

import (
	"context"

	"github.com/m-mizutani/transmith/pkg/domain/model"
)

type UseCases interface {
	Transmit(ctx context.Context, msg model.Message) error
}
