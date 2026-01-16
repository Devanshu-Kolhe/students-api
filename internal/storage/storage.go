package storage

import (
	// "github.com/Devanshu-Kolhe/students-api/internal/http/handlers/student"
	"github.com/Devanshu-Kolhe/students-api/internal/types"
)

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateStudent(student types.Student) (types.Student, error)
}
