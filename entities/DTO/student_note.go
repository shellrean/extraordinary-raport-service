package dto

type StudentNoteRequest struct {
	ID			int64 		`json:"id"`
	Type 		int 		`json:"type" validate="required"`
	TeacherID 	int64 		`json:"teacherID"`
	StudentID 	int64 		`json:"studentID" validate="required"`
	Note 		string 		`json:"note" validate="required"`
}

type StudentNoteResponse struct {
	ID			int64 		`json:"id"`
	Type 		int 		`json:"type"`
	TeacherID 	int64 		`json:"teacherID"`
	StudentID 	int64 		`json:"studentID"`
	Note 		string 		`json:"note"`
}