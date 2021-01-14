package domain

import (
	"time"
)

// Student Result Plan
type StudentResultPlan struct {
	ID 			int64
	Student 	ClassroomStudent
	Subject 	ClassroomSubject
	Plan		ClassroomSubjectPlan
	Number 		uint
	UpdatedBy	User
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}