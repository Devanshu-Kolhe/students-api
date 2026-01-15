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
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
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
