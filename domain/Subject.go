package domain

import (
	"time"
	"context"
)

type Subject struct {
	ID 			int64
	Name 		string
	Type 		string
	CreatedAt	time.Time
	UpdatedAt 	time.Time
}

type SubjectRepository interface {
	Fetch(ctx context.Context, cursor int64, num int64) ([]Subject, error)
}

type SubjectUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Subject, string, error)
}