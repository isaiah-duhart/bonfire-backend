package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJson(w, code, ErrorResponse{
		Error: msg,
	})
}

func RespondWithJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	json, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling json: ", err)
		return
	}

	_, err = w.Write(json)
	if err != nil {
		fmt.Println("Error writing json: ", err)
		return
	}
}