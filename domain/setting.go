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
	GetByName(ctx context.Context, name string) (Setting, error)
}