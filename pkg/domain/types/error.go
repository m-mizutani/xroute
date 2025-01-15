package types

import "github.com/m-mizutani/goerr/v2"

var (
	ErrTagUnauthorized = goerr.NewTag("unauthorized")
	ErrTagBadRequest   = goerr.NewTag("bad_request")
)
