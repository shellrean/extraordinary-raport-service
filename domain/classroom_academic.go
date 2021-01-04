package domain

import (
	"context"
	"time"
)

type ClassroomAcademic struct {
	ID 			int64		`json:"id"`
	Academic 	Academic	`json:"academic"`
	Classroom 	Classroom	`json:"classroom"`
	Teacher 	User		`json:"teacher"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}

type ClassroomAcademicRepository interface {
	Fetch(ctx context.Context) ([]ClassroomAcademic, error)
	Store(ctx context.Context, ca *ClassroomAcademic) (error)
	GetByAcademicAndClass(ctx context.Context, a int64, c int64) (ClassroomAcademic, error)
}

type ClassroomAcademicUsecase interface {
	Fetch(ctx context.Context) ([]ClassroomAcademic, error)
	Store(ctx context.Context, ca *ClassroomAcademic) (error)
}