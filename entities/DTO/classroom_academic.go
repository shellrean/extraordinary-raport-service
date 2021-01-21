package dto

type ClassroomAcademicResponse struct {
	ID 				int64	`json:"id"`
	AcademicID 		int64	`json:"academicID"`
	TeacherID 		int64	`json:"teacherID"`
	ClassroomID 	int64	`json:"classroomID"`
	TeacherName 	string	`json:"teacherName"`
	ClassroomName 	string 	`json:"classroomName"`
}

type ClassroomAcademicRequest struct {
	ID 			int64		`json:"id"`
	TeacherID 	int64		`json:"teacherID" validate:"required"`
	ClassroomID int64 		`json:"classroomID" validate:"required"`
}