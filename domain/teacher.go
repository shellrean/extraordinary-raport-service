package domain

import (
	"time"
	// "context"
)

type Teacher struct {
	ID 			int64 		`json:"id"`
	NIP 		string 		`json:"nip"`
	CreatedAt   time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
}

type TeacherRepository interface {

} 

type TeacherUsecase interface {
	
}