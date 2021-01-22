package dto

type ClassroomSubjectResponse struct {
	ID 					int64		`json:"id"`
	ClassroomAcademicID int64		`json:"classroomAcademicID"`
	SubjectID 			int64		`json:"subjectID"`
	SubjectName 		string 		`json:"subjectName"`
	TeacherID 			int64		`json:"teacherID"`
	TeacherName 		string 		`json:"teacherName"`
	MGN					int 		`json:"mgn"`
}