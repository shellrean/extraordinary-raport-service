package domain

import (
	"context"
	"time"
)

type ClassroomAcademic struct {
	ID 			int64		`json:"id"`
	Academic 	Academic	`json:"academic"`
	Teacher 	User		`json:"teacher"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}

type ClassroomAcademicRepository interface {
	Fetch(ctx context.Context) ([]ClassroomAcademic, error)
}

type ClassroomAcademicUsecase interface {
	Fetch(ctx context.Context) ([]ClassroomAcademic, error)
}