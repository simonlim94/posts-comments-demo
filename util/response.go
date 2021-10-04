package util

import (
	"encoding/json"
	"net/http"
)

type ResponseBody struct {
	StatusCode int         `json:"statusCode,omitempty"`
	Body       string      `json:"body,omitempty"`
	Items      interface{} `json:"items,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}

var (
	InternalServerErrorResponseBody = ResponseBody{
		StatusCode: http.StatusInternalServerError,
		Body:       "An internal server occured",
	}
	NotFoundErrorResponseBody = ResponseBody{
		StatusCode: http.StatusNotFound,
		Body:       "This API endpoint is not found",
	}
	BadRequestErrorResponseBody = ResponseBody{
		StatusCode: http.StatusBadRequest,
		Body:       "Invalid request body is provided",
	}
)

func WriteOKResponse(w http.ResponseWriter, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := &ResponseBody{}
	resp.StatusCode = http.StatusOK
	resp.Items = payload

	return json.NewEncoder(w).Encode(resp)
}
