package dto

type ClassroomSubjectPlanRequest struct {
	ID 			int64 		`json:"id"`
	Type 		string 		`json:"type" validate:"required"`
	Name 		string 		`json:"name" validate:"required"`
	Desc 		string 		`json:"desc" validate:"required"`
	TeacherID 	int64 		`json:"teacherID" validate:"required"`
	SubjectID 	int64 		`json:"classroomSubjectID" validate:"required"`
	ClassroomID int64 		`json:"classroomAcademicID" validate:"required"`
	CountPlan 	uint 		`json:"countPlan" validate:"required"`
	MaxPoint 	uint 		`json:"maxPoint" validate:"required"`
}