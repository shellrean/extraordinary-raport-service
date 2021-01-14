package domain

import (
	"time"
)

type StudentResultExschool struct {
	ID 					int64
	ExschoolStudent 	ExschoolStudent
	Number				uint
	UpdatedBy			User
	CreatedAt			time.Time
	UpdatedAt 			time.Time
}