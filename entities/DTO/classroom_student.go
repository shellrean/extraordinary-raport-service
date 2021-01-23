package dto

type ClassroomStudentResponse struct {
	ID 					int64		`json:"id"`
	ClassroomAcademicID int64		`json:"classroomAcademicID"`
	StudentID 			int64		`json:"studentID"`
	StudentSRN 			string 		`json:"studentSRN"`
	StudentNSRN			string 		`json:"studentNSRN"`
	StudentName 		string 		`json:"studentName"`
}

type ClassroomStudentRequest struct {
	ID 					int64		`json:"id"`
	ClassroomAcademicID int64		`json:"classroomAcademicID" validate:"required"`
	StudentID 			int64		`json:"studentID" validate:"required"`
}

type ClassroomStudentCopyRequest struct {
	ClassroomAcademicID 	int64 		`json:"classroomAcademicID"`
	ToClassroomAcademicID 	int64 		`json:"toClassroomAcademicID"`
}