package domain

import (
	"time"
	"context"
)

// Classroom Subject Plan Result
type ClassroomSubjectPlanResult struct {
	ID 			int64
	Student 	ClassroomStudent
	Subject 	ClassroomSubject
	Plan		ClassroomSubjectPlan
	Number 		uint
	UpdatedBy	User
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}

type ClassroomSubjectPlanResultRepository interface {
	Store(ctx context.Context, spr *ClassroomSubjectPlanResult) (error)
}

type ClassroomSubjectPlanResultUsecase interface {
	Store(ctx context.Context, spr *ClassroomSubjectPlanResult) (error)
}