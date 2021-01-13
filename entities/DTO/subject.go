package dto

type SubjectResponse struct {
	ID 		int64 		`json:"id"`
	Name 	string 		`json:"name" validate:"required"`
	Type 	string 		`json:"type" validate:"required"`
}