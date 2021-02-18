package domain

import (
	"time"
	"context"
)

type Exschool struct {
	ID 			int64
	Name 		string
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}

type ExschoolRepository interface {
	Fetch(ctx context.Context) ([]Exschool, error)
	GetByID(ctx context.Context, id int64) (Exschool, error)
	Store(ctx context.Context, ex *Exschool) (error)
}

type ExschoolUsecase interface {
	Fetch(ctx context.Context) ([]Exschool, error)
	GetByID(ctx context.Context, id int64) (Exschool, error)
	Store(ctx context.Context, ex *Exschool) (error)
}