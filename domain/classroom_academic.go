package domain

import (
	"context"
	"time"
)

type ClassroomAcademic struct {
	ID 			int64
	Academic 	Academic
	Classroom 	Classroom
	Teacher 	User
	CreatedAt	time.Time
	UpdatedAt 	time.Time
}

type ClassroomAcademicRepository interface {
	Fetch(ctx context.Context, academicID int64) ([]ClassroomAcademic, error)
	Store(ctx context.Context, ca *ClassroomAcademic) (error)
	GetByID(ctx context.Context, id int64) (ClassroomAcademic, error)
	GetByAcademicAndClass(ctx context.Context, a int64, c int64) (ClassroomAcademic, error)
	Update(ctx context.Context, ca *ClassroomAcademic) (error)
	Delete(ctx context.Context, id int64) (error)
}

type ClassroomAcademicUsecase interface {
	Fetch(ctx context.Context) ([]ClassroomAcademic, error)
	FetchByAcademic(ctx context.Context, academicID int64) ([]ClassroomAcademic, error)
	GetByID(ctx context.Context, id int64) (ClassroomAcademic, error)
	Store(ctx context.Context, ca *ClassroomAcademic) (error)
	Update(ctx context.Context, ca *ClassroomAcademic) (error)
	Delete(ctx context.Context, id int64) (error)
}