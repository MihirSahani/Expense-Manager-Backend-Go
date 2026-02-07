package ehandler

import (
	"net/http"
)

func Health(version string, environment string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: In a production system, you would likely log the error from WriteJSON.
		_ = WriteJSON(w, http.StatusOK, map[string]string{
			"status":      "ok",
			"version":     version,
			"environment": environment,
		})
	})
}
