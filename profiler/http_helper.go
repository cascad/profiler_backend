package profiler

import (
	"net/http"
	"encoding/json"
)

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithStream(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJson(w, code, map[string]string{"error": msg})
}
