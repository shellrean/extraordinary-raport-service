package domain

import (
	"time"
	"context"
)

type Classroom struct {
	ID 			int64 		`json:"id"`
	Name 		string 		`json:"name"`
	Grade 		string 		`json:"grade"`
	Major 		Major 		`json:"major"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}

type ClassroomRepository interface {
	Fetch(ctx context.Context) ([]Classroom, error)
	GetByID(ctx context.Context, id int64) (Classroom, error)
	Store(ctx context.Context, c *Classroom) (error)
	Update(ctx context.Context, c *Classroom) (error)
	Delete(ctx context.Context, id int64) (error)
}

type ClassroomUsecase interface {
	Fetch(ctx context.Context) ([]Classroom, error)
	GetByID(ctx context.Context, id int64) (Classroom, error)
	Store(ctx context.Context, c *Classroom) (error)
	Update(ctx context.Context, c *Classroom) (error)
	Delete(ctx context.Context, id int64) (error)
}