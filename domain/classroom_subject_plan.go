package domain

import (
	"time"
)

const (
	PlanTask = "task" // Type for task
	PlanEVMS = "evms" // Type for evaluation middle semester
	PlanEVLS = "evls" // Type for evaluation last semester
	PlanExam = "exam" // Type for exam
)

type ClassroomSubjectPlan struct {
	ID 			int64
	Type 		string
	Name 		string
	Desc 		string
	Teacher 	User
	Classroom 	ClassroomAcademic
	CountPlan 	uint
	MaxPoint 	uint
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}