package response

import (
	"encoding/json"
	"net/http"
)

type (
	Response struct {
		Status  int         `json:"status,omitempty"`
		Message string      `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}
)

func WriteSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	resp := &Response{
		Status:  statusCode,
		Message: "Success",
		Data:    data,
	}

	writeResp(w, statusCode, resp)
}

func WriteError(w http.ResponseWriter, statusCode int, errMessage string) {
	resp := &Response{
		Status:  statusCode,
		Message: errMessage,
	}

	writeResp(w, statusCode, resp)
}

func writeResp(w http.ResponseWriter, statusCode int, resp *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
