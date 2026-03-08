package jsonutil

import (
	"encoding/json"
	"net/http"
)

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
