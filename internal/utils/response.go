package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error string `json:"error,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

const (
	StatusOK = "OK"
	StatusError = "ERROR"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error{

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {

	return Response{
		Status: StatusError,
		Error: err.Error(),
		Data: nil,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errorsMsg []string
	for _, err := range errs {

		switch err.Tag() {
		case "required":
			errorsMsg = append(errorsMsg, fmt.Sprintf("%s is required", err.Field()))
		default:
			errorsMsg = append(errorsMsg, err.Error())
		}
		
	}
	
	return Response{
		Status: StatusError,
		Error: strings.Join(errorsMsg, ", "),
	}
}