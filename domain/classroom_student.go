package domain

import (
	"time"
	"context"
)

type ClassroomStudent struct {
	ID 					int64
	ClassroomAcademic 	ClassroomAcademic
	Student 			Student
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
}

type ClassroomStudentRepository interface {
	Fetch(ctx context.Context, cursor int64, num int64) ([]ClassroomStudent, error)
	GetByID(ctx context.Context, id int64) (ClassroomStudent, error)
	Store(ctx context.Context, cs *ClassroomStudent) (error)
	Update(ctx context.Context, cs *ClassroomStudent) (error)
	Delete(ctx context.Context, id int64) (error)
}

type ClassroomStudentUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]ClassroomStudent, string, error)
	GetByID(ctx context.Context, id int64) (ClassroomStudent, error)
	Store(ctx context.Context, cs *ClassroomStudent) (error)
	Update(ctx context.Context, cs *ClassroomStudent) (error)
	Delete(ctx context.Context, id int64) (error)
}