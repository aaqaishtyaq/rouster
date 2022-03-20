package builder

import "context"

type Builder interface {
	Build(ctx *context.Context)
}

type ErrorLine struct {
	Error        string      `json:"error"`
	ErrorDetails ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}
