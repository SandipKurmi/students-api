package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/SandipKurmi/students-api/internal/types"
	response "github.com/SandipKurmi/students-api/internal/utils"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {


		var student  types.Student

	  	err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF){
			response.WriteJson(w, http.StatusAlreadyReported, response.GeneralError(fmt.Errorf("empty request body")) )

			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}


		if err := validator.New().Struct(student); err != nil {

			errs := err.(validator.ValidationErrors)

			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(errs))
			return
		}

		


		response.WriteJson(w, http.StatusCreated, map[string]string{"success":"OK"})
	}
}