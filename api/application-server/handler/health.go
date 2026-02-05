package handler

import (
	"net/http"
	"github.com/krakn/expense-management-backend-go/internal/utils"
)

func Health(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{
		"version": "0.0.1",
		"environment": utils.GetEnv("ENVIRONMENT", "unknown"),
	})
}