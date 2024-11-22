package routes

import (
	"net/http"
	"paws/internal/response"
)

type PingPongHandler struct{}

func NewPingPongHandler() *PingPongHandler {
	return &PingPongHandler{}
}

func (h *PingPongHandler) RegisterRoutes(mux *http.ServeMux, mf MiddlewareFunc) {
	mux.HandleFunc("GET /api/v1/ping", h.Ping)
}

func (h *PingPongHandler) Ping(w http.ResponseWriter, r *http.Request) {
	response.Text(w, "pong")
}
