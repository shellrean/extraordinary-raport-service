package domain

import (
	"time"
)

type ExschoolStudent struct {
	ID 			int64
	Exschool 	Exschool
	Student 	ClassroomStudent
	UpdatedBy	User
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}