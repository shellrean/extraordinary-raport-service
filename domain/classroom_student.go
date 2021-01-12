package domain

import (
	"time"
)

type ClassroomStudent struct {
	ID 					int64
	Academic 			Academic
	ClassroomAcademic 	ClassroomAcademic
	Student 			Student
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
}