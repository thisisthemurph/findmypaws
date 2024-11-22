package response

import (
	"encoding/json"
	"io"
	"net/http"
)

type apiResponse struct {
	writer http.ResponseWriter
}

func WithStatus(w http.ResponseWriter, status int) apiResponse {
	w.WriteHeader(status)
	return apiResponse{w}
}

func (r apiResponse) SendJSON(v any) {
	JSON(r.writer, v)
}

func JSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error encoding response"})
	}
}

func Text(w http.ResponseWriter, v string) {
	w.Header().Set("Content-Type", "text/plain")
	if _, err := io.WriteString(w, v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error encoding response"})
	}
}
