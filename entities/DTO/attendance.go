package dto

type AttendanceResponse struct {
	ID			int64 		`json:"id"`
	StudentID	int64		`json:"studentID"`
	Type 		int 		`json:"type"`
	Total 		uint8 		`json:"total"`
}

type AttendanceRequest struct {
	ID 			int64		`json:"id"`
	StudentID 	int64		`json:"studentID" validate:"required"`
	Type 		int			`json:"type" validate:"required"`
	Total 		uint8 		`json:"total" validate:"required"`
}