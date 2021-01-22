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
	Store(ctx context.Context, cs *ClassroomSubject) (error)
	GetByID(ctx context.Context, id int64) (ClassroomSubject, error)
	GetByClassroomAndSubject(ctx context.Context, academicClassroomID int64, subjectID int64) (ClassroomSubject, error)
	Update(ctx context.Context, cs *ClassroomSubject) (error)
	Delete(ctx context.Context, idi int64) (error)
}

type ClassroomSubjectUsecase interface {
	FetchByClassroom(ctx context.Context, academicClassroomID int64) ([]ClassroomSubject, error)
	Store(ctx context.Context, cs *ClassroomSubject) (error)
	GetByID(ctx context.Context, id int64) (ClassroomSubject, error)
	Update(ctx context.Context, cs *ClassroomSubject) (error)
	Delete(ctx context.Context, idi int64) (error)
}