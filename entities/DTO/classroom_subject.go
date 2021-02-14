package dto

type ClassroomSubjectResponse struct {
	ID 					int64		`json:"id"`
	ClassroomAcademicID int64		`json:"classroomAcademicID"`
	ClassroomName 		string 		`json:"classroomName"`
	SubjectID 			int64		`json:"subjectID"`
	SubjectName 		string 		`json:"subjectName"`
	TeacherID 			int64		`json:"teacherID"`
	TeacherName 		string 		`json:"teacherName"`
	MGN					int 		`json:"mgn"`
}

type ClassroomSubjectRequest struct {
	ID 					int64		`json:"id"`
	ClassroomAcademicID int64		`json:"classroomAcademicID" validate:"required"`
	SubjectID 			int64		`json:"subjectID" validate:"required"`
	TeacherID			int64 		`json:"teacherID", validate:"required"`		
	MGN					int 		`json:"mgn"`
}

type ClassroomSubjectCopyRequest struct {
	ClassroomAcademicID 	int64 	`json:"classroomAcademicID" validate:"required"`
	ToClassroomAcademicID	int64 	`json:"toClassroomAcademicID" validate:"required"`
}