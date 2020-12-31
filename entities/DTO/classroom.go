package dto

type ClassroomResponse struct {
	ID 			int64 		`json:"id"`
	Name 		string 		`json:"name"`
	Grade 		string 		`json:"grade"`
	MajorID 	int64 		`json:"major_id"`
}