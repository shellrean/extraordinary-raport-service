package domain

import (
	"time"
	"context"
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
	Subject 	ClassroomSubject
	Classroom 	ClassroomAcademic
	CountPlan 	uint
	MaxPoint 	uint
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}

type ClassroomSubjectPlanRepository interface {
	GetByID(ctx context.Context, id int64) (ClassroomSubjectPlan, error)
	FetchByClassroom(ctx context.Context, id int64) ([]ClassroomSubjectPlan, error)
	FetchByTeacher(ctx context.Context, id int64) ([]ClassroomSubjectPlan, error)
	FetchByTeacherAndClassroom(ctx context.Context, tid int64, cid int64) ([]ClassroomSubjectPlan, error)
	Store(ctx context.Context, csp *ClassroomSubjectPlan) (error)
	Update(ctx context.Context, csp *ClassroomSubjectPlan) (error)
	Delete(ctx context.Context, id int64) (error)
}

type ClassroomSubjectPlanUsecase interface {
	Fetch(ctx context.Context, query string, userID int64, classID int64) ([]ClassroomSubjectPlan, error)
	GetByID(ctx context.Context, id int64) (ClassroomSubjectPlan, error)
	Store(ctx context.Context, csp *ClassroomSubjectPlan) (error)
	Update(ctx context.Context, csp *ClassroomSubjectPlan) (error)
	Delete(ctx context.Context, id int64) (error)
}