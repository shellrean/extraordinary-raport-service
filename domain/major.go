package domain

import (
	"time"
	"context"
)

type Major struct {
	ID 			int64 		`json:"id"`
	Name 		string 		`json:"name"`
	CreatedAt	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}

type MajorRepository interface {
	Fetch(ctx context.Context) ([]Major, error)
	GetByID(ctx context.Context, id int64) (Major, error)
}

type MajorUsecase interface {
	Fetch(ctx context.Context) ([]Major, error)
}