package domain

import (
	"context"
	"time"
)

const (
	NoteDailly = 1
	NoteGrd = 2
	NoteAcademic = 3
	NoteCharIntegrity = 4
	NoteCharReligius = 5
	NoteCharNation = 6
	NoteCharIndependence = 7
	NoteCharTeamwork = 8
)

type StudentNote struct {
	ID 			int64
	Type 		int
	Teacher 	User
	Student 	ClassroomStudent
	Note 		string
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}

type StudentNoteRepository interface {
	Store(ctx context.Context, sn *StudentNote) (error)
	FetchByClassroom(ctx context.Context, id int64) ([]StudentNote, error)
	FetchByTypeAndClassroom(ctx context.Context, id int64, typ int64) ([]StudentNote, error)
	GetByStudentAndType(ctx context.Context, id int64, typ int) (StudentNote, error)
	Update(ctx context.Context, sn *StudentNote) (error)
}

type StudentNoteUsecase interface {
	Store(ctx context.Context, sn *StudentNote) (error)
	FetchByClassroom(ctx context.Context, id int64, typ int64) ([]StudentNote, error)
}