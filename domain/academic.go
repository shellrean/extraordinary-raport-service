package domain

import (
	"time"
	"context"
)

type Academic struct {
	ID 			int64 		`json:"id"`
	Name 		string 		`json:"name"`
	Semester 	uint8 		`json:"semester"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}

type AcademicRepository interface {
	Fetch(ctx context.Context) ([]Academic, error)
	GetByID(ctx context.Context, id int64) (Academic, error)
	GetByYearAndSemester(ctx context.Context, year string, semester int) (Academic, error)
	Store(ctx context.Context, ac *Academic) (error)
	Delete(ctx context.Context, id int64) (error)
}

type AcademicUsecase interface {
	Fetch(ctx context.Context) ([]Academic, error)
	Generate(ctx context.Context) (Academic, error)
}