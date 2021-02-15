package dto

type ClassroomSubjectPlanResultRequest struct {
	ID 			int64 	`json:"id"`
	Index 		int		`json:"index"`
	StudentID	int64	`json:"studentID"`
	SubjectID 	int64 	`json:"subjectID"`
	PlanID 		int64 	`json:"planID"`
	Number 		uint 	`json:"number"`
}

type ClassroomSubjectPlanResultResponse struct {
	ID 			int64  	`json:"id"`
	Index 		int 	`json:"index"`
	StudentID 	int64 	`json:"studentID"`
	SubjectID 	int64 	`json:"subjectID"`
	PlanID 		int64 	`json:"planID"`
	Number 		uint 	`json:"number"`
	UpdatedByID int64 	`json:"updatedByID"`
}