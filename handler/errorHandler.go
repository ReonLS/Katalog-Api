package handler

import (
	"encoding/json"
	"net/http"
	"simple-product-api/models"
)

func GenerateError(rw http.ResponseWriter, message string, status_code int ){
	//Content Type Application/problem+json for production, for now no need
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(status_code)

	errResp := &models.ErrorResponse{
		Message: message,
		StatusCode: status_code,
	}

	json.NewEncoder(rw).Encode(errResp)
}