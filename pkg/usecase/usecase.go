package usecase

import "github.com/m-mizutani/transmith/pkg/adapter"

type UseCases struct {
	adaptors *adapter.Adapters
}

func New(adaptors *adapter.Adapters) *UseCases {
	return &UseCases{adaptors: adaptors}
}
