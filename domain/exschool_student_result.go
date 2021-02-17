package domain

import (
	"time"
)

type ExschoolStudentResult struct {
	ID 					int64
	ExschoolStudent 	ExschoolStudent
	Number				string
	UpdatedBy			User
	CreatedAt			time.Time
	UpdatedAt 			time.Time
}