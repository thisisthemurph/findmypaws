package routes

import (
	"encoding/json"
	"net/http"
)

type PingPongHandler struct{}

func NewPingPongHandler() *PingPongHandler {
	return &PingPongHandler{}
}

func (h *PingPongHandler) RegisterRoutes(mux *http.ServeMux, mf MiddlewareFunc) {
	mux.HandleFunc("GET /api/v1/ping", h.Ping)
}

func (h *PingPongHandler) Ping(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("pong")
}
