package handler

import (
	"api-gateway/internal/deps"
)

type Handler struct {
	deps *deps.Dependencies
}

func NewHandler(dependencies *deps.Dependencies) *Handler {
	return &Handler{
		deps: dependencies,
	}
}
