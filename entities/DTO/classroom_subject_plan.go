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

type ClassroomSubjectPlanResponse struct {
	ID 			int64 		`json:"id"`
	Type 		string 		`json:"type"`
	Name 		string 		`json:"name"`
	Desc 		string 		`json:"desc"`
	TeacherID 	int64 		`json:"teacherID"`
	SubjectID 	int64 		`json:"classroomSubjectID"`
	SubjectName string		`json:"subjectName"`
	ClassroomID int64 		`json:"classroomAcademicID"`
	CountPlan 	uint 		`json:"countPlan"`
	MaxPoint 	uint 		`json:"maxPoint"`
}

type ClassroomSubjectPlanFetchRequest struct {
    ClassroomID int64       `json:"classroomID"`
    TeacherID   int64       `json:"teacherID"`
    Query       string      `json:"query"`
}