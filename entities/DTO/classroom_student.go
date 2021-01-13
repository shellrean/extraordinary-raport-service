package dto

type ClassroomStudentResponse struct {
	ID 					int64		`json:"id"`
	ClassroomAcademicID int64		`json:"classroom_academic_id" validate:"required"`
	StudentID 			int64		`json:"student_id" validate:"required"`
}