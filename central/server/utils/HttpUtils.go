package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func ErrorResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	b := &ErrorBody{
		Status:  "FAIL",
		Message: message,
	}
	msg, _ := json.Marshal(b)
	_, _ = w.Write(msg)
}

type SuccessBody struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponseWithData(w http.ResponseWriter, code int, message string, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	b := &SuccessBody{
		Status:  "SUCCESS",
		Message: message,
		Data:    data,
	}
	msg, _ := json.Marshal(b)
	_, _ = w.Write(msg)
}

func SuccessResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	b := &SuccessBody{
		Status:  "SUCCESS",
		Message: message,
	}
	msg, _ := json.Marshal(b)
	_, _ = w.Write(msg)
}
