package domain

import (
	"time"
	"context"
)

type ExschoolStudent struct {
	ID 			int64
	Exschool 	Exschool
	Student 	ClassroomStudent
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}

type ExschoolStudentRepository interface {
	FetchByClassroom(ctx context.Context, classroomID int64) ([]ExschoolStudent, error)
	GetByID(ctx context.Context, id int64) (ExschoolStudent, error)
	Store(ctx context.Context, es *ExschoolStudent) (error)
	Delete(ctx context.Context, id int64) (error)
}

type ExschoolStudentUsecase interface {
	FetchByClassroom(ctx context.Context, classroomID int64) ([]ExschoolStudent, error)
	GetByID(ctx context.Context, id int64) (ExschoolStudent, error)
	Store(ctx context.Context, es *ExschoolStudent) (error)
	Delete(ctx context.Context, id int64) (error)
}