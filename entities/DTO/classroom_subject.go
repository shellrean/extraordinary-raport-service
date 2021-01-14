package dto

type ClassroomSubjectResponse struct {
	ID 					int64		`json:"id"`
	ClassroomAcademicID int64		`json:"classroom_academic_id"`
	SubjectID 			int64		`json:"subject_id"`
	TeacherID 			int64		`json:"teacher_id"`
	MGN					int 		`json:"mgn"`
}