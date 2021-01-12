package domain

import (
	"context"
	"time"
)

type Parent struct {
	Name 			string
	Address 		string
	Profession		string
	Telp			string
}

type Familly struct {
	Status			string
	Order			string
}

type Student struct {
	ID				int64
	SRN				string
	NSRN			string
	Name			string 
	Gender			string
	BirthPlace		string
	BirthDate		string
	Religion		Religion
	Address			string
	Telp			string
	SchoolBefore	string
	AcceptedGrade	string
	AcceptedDate	string
	Familly			Familly
	Father 			Parent
	Mother			Parent
	Guardian		Parent
	CreatedAt		time.Time
	UpdatedAt		time.Time
}

type StudentRepository interface {
	Fetch(ctx context.Context, cursor int64, num int64) ([]Student, error)
	GetByID(ctx context.Context, id int64) (Student, error)
	Store(ctx context.Context, s *Student) (error)
	Update(ctx context.Context, s *Student) (error)
}

type StudentUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Student, string, error)
	GetByID(ctx context.Context, id int64) (Student, error)
	Store(ctx context.Context, s *Student) (error)
	Update(ctx context.Context, s *Student) (error)
}