package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Devanshu-Kolhe/students-api/internal/types"
	"github.com/Devanshu-Kolhe/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func Create() http.HandlerFunc {
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

		response.Writejson(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}
