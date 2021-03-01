package dto

type StudentNoteRequest struct {
	ID			int64 		`json:"id"`
	Type 		int 		`json:"type" validate="required"`
	TeacherID	int64 		`json:"teacherID" validate="required"`
	StudentID 	int64 		`json:"studentID" validate="required"`
	Note 		string 		`json:"note" validate="required"`
}