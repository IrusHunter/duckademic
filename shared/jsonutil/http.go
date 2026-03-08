package jsonutil

import (
	"encoding/json"
	"net/http"
)

// ResponseWithJSON writes the given data as a JSON response with the specified HTTP status.
//
// It requires an http.ResponseWriter (w), a status code (status), and any value that can be marshaled to JSON (data).
//
// If marshaling fails, it responds with HTTP 500 Internal Server Error.
func ResponseWithJSON(w http.ResponseWriter, status int, data any) {
	dat, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(dat)
}

func ResponseWithError(w http.ResponseWriter, status int, err error) {
	ResponseWithJSON(w, status, map[string]string{"error": err.Error()})
}
