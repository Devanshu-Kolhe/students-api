package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Devanshu-Kolhe/students-api/internal/storage"
	"github.com/Devanshu-Kolhe/students-api/internal/types"
	"github.com/Devanshu-Kolhe/students-api/internal/utils/response"
	// "github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func Create(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Creating a student")
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.Writejson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return

		}

		if err != nil {
			response.Writejson(w, http.StatusBadGateway, response.GeneralError(err))
			return
		}

		//request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.Writejson(w, http.StatusBadGateway, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("User Created Succeessfully", slog.String("id", fmt.Sprint(lastId)))
		if err != nil {
			response.Writejson(w, http.StatusInternalServerError, err)
			return
		}
		response.Writejson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a Student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.Writejson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.GetStudentById(intId)

		if err != nil {
			slog.Error("Error getting user")
			response.Writejson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.Writejson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")

		students, err := storage.GetStudents()
		if err != nil {
			response.Writejson(w, http.StatusInternalServerError, err)
			return
		}

		response.Writejson(w, http.StatusOK, students)
	}
}

func UpdateStudentById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var student types.Student
		slog.Info("Updating a Student..")

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		idStr := r.PathValue("id")
		if idStr == "" {
			response.Writejson(w, http.StatusBadRequest,
				response.GeneralError(fmt.Errorf("student id is required")))
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			response.Writejson(w, http.StatusBadRequest,
				response.GeneralError(fmt.Errorf("invalid student id")))
			return
		}

		if err := decoder.Decode(&student); err != nil {
			if errors.Is(err, io.EOF) {
				response.Writejson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty request body")))
				return
			}
			response.Writejson(w, http.StatusBadRequest,
				response.GeneralError(err))
			return
		}
		if student.Id != 0 && student.Id != id {
			response.Writejson(w, http.StatusBadRequest,
				response.GeneralError(fmt.Errorf("id in URL and body do not match")))
			return
		}

		student.Id = id

		var validate = validator.New()
		if err := validate.Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.Writejson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		updatedStudent, err := storage.UpdateStudent(student)
		if err != nil {
			slog.Error("Failed to Update student", 
			"id",id,	
			"error", err,
		)
			response.Writejson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info("Student Updated successfully", "id", student.Id)

		response.Writejson(w, http.StatusOK, updatedStudent)
	}
}
