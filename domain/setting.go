package domain

import (
	"time"
	"context"
)

const (
	SettingAcademicActive 	= "academic_active"
)

type Setting struct {
	ID 			int64 		`json:"id"`
	Name 		string 		`json:"name"`
	Value 		string 		`json:"value"`
	CreatedAt 	time.Time	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}

type SettingRepository interface {
	Fetch(ctx context.Context, names []string) ([]Setting, error)
	GetByName(ctx context.Context, name string) (Setting, error)
	Update(ctx context.Context, s *Setting) (error)
}

type SettingUsecase interface {
	Fetch(ctx context.Context, names []string) ([]Setting, error)
	Update(ctx context.Context, s *Setting) (error)
}