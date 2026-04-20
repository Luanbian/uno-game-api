// Package health provides a HTTP GET/ route to verify server status
package health

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

const Route = "GET /health"

type DataType struct {
	Message string
	Version string
}

type APIResponse struct {
	TransactionID string
	Message       string
	Data          DataType
	Error         *string
}

func Handler(res http.ResponseWriter, req *http.Request) {
	responseData := DataType{
		Message: "API is healthy",
		Version: "0.0.1",
	}

	response := APIResponse{
		TransactionID: uuid.NewString(),
		Message:       "Success",
		Error:         nil,
		Data:          responseData,
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(res).Encode(response); err != nil {
		http.Error(res, "Erro ao escrever resposta", http.StatusInternalServerError)
	}
}
