package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/SandipKurmi/students-api/internal/storage"
	"github.com/SandipKurmi/students-api/internal/types"
	response "github.com/SandipKurmi/students-api/internal/utils"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
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

	    id, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		slog.Info("Student created", slog.String("id", fmt.Sprintf("%d", id)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"id": fmt.Sprintf("%d", id), "status":"OK"})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Get the student ID from the request URL
		id := r.PathValue("id")

		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("id is required")))
			return
		}

		// Convert the ID to an integer
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id")))
			return
		}

		// Get the student from the storage
		student, err := storage.GetStudentById(idInt)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		// Write the student as JSON
		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetAll(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get all students from the storage
		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		// Write the students as JSON
		response.WriteJson(w, http.StatusOK, students)
	}
}

func Update(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the student ID from the request URL
		id := r.PathValue("id")
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("id is required")))
			return
		}
		// Convert the ID to an integer
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id")))
			return
		}
		// Get the student from the request body
		var student types.Student
		err = json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		// Update the student in the storage
		_, err = storage.UpdateStudent(idInt, student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		// Write a success response
		response.WriteJson(w, http.StatusOK, map[string]string{"status": "OK"})
	}
}

func Delete(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the student ID from the request URL
		id := r.PathValue("id")
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("id is required")))
			return
		}
		// Convert the ID to an integer
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id")))
			return
		}
		// Delete the student from the storage
		_, err = storage.DeleteStudent(idInt)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		// Write a success response
		response.WriteJson(w, http.StatusOK, map[string]string{"status": "OK"})
	}
}
