package domain

import (
	"time"
	"context"
)

type ClassroomSubject struct {
	ID					int64
	ClassroomAcademic	ClassroomAcademic
	Subject 			Subject
	Teacher				User
	MGN					int
	CreatedAt 			time.Time
	UpdatedAt 			time.Time
}

type ClassroomSubjectRepository interface {
	FetchByClassroom(ctx context.Context, academicClassroomID int64) ([]ClassroomSubject, error)
}

type ClassroomSubjectUsecase interface {
	FetchByClassroom(ctx context.Context, academicClassroomID int64) ([]ClassroomSubject, error)
}