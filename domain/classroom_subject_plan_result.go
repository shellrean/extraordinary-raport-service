package domain

import (
	"time"
	"context"
)

// Classroom Subject Plan Result
type ClassroomSubjectPlanResult struct {
	ID 			int64
	Index 		int
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
	FetchByPlan(ctx context.Context, planID int64) ([]ClassroomSubjectPlanResult, error)
	Update(ctx context.Context, spr *ClassroomSubjectPlanResult) (error)
	GetByPlanIndexStudent(ctx context.Context, planID int64, idx int, studentID int64) (ClassroomSubjectPlanResult, error)
}

type ClassroomSubjectPlanResultUsecase interface {
	Store(ctx context.Context, spr *ClassroomSubjectPlanResult) (error)
	FetchByPlan(ctx context.Context, planID int64) ([]ClassroomSubjectPlanResult, error)
}