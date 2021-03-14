package domain

import (
	"context"
	"time"
)

const (
	TypeAttendSick = 1
	TypeAttendPermit = 2
	TypeAttendNull = 3
)

type Attendance struct {
	ID 			int64
	Student		ClassroomStudent
	Type 		int
	Total 		uint8
	CreatedAt	time.Time
	UpdatedAt 	time.Time
}

type AttendanceRepository interface {
	Fetch(ctx context.Context, cid int64) ([]Attendance, error)
	Store(ctx context.Context, a *Attendance) (error)
	Update(ctx context.Context, a *Attendance) (error)
	GetByStudentAndType(ctx context.Context, sid int64, typ int) (Attendance, error)
}

type AttendanceUsecase interface {
	Fetch(ctx context.Context, cid int64) ([]Attendance, error)
	Store(ctx context.Context, a *Attendance) (error)
}