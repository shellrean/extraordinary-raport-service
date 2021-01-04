package dto

type ClassroomAcademicResponse struct {
	ID 			int64		`json:"id"`
	AcademicID 	int64		`json:"academic_id" validate:"required"`
	TeacherID 	int64		`json:"teacher_id" validate:"required"`
	ClassroomID int64 		`json:"classroom_id" validate:"required"`
}