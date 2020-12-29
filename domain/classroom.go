package domain

import (
	"time"
	// "context"
)

type Classroom struct {
	ID 			int64 		`json:"id"`
	Name 		string 		`json:"name"`
	Grade 		string 		`json:"grade"`
	Major 		Major 		`json:"major"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}