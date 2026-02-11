package ehandler

import (
	"encoding/json"
	"net/http"
)

func ReadJSON(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(&data)
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	output, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(output)

	return err 
}