package dto

type ClassroomAcademicResponse struct {
	ID 			int64		`json:"id"`
	AcademicID 	int64		`json:"academic_id"`
	TeacherID 	int64		`json:"teacher_id"`
}