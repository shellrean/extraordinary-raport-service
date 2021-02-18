package dto

type ExschoolStudentResponse struct {
	ID 			int64 		`json:"id"`
	ExschoolID	int64 		`json:"exschoolID"`
	StudentID 	int64 		`json:"studentID"`
}

type ExschoolStudentRequest struct {
	ID 			int64 		`json:"id"`
	ExschoolID	int64 		`json:"exschoolID" valiate:"required"`
	StudentID 	int64 		`json:"studentID" validate:"required"`
}