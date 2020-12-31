package dto

type AcademicResponse struct {
	ID 			int64 		`json:"id"`
	Name 		string 		`json:"name"`
	Semester 	uint8 		`json:"semester"`	
}