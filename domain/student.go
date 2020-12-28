package domain

import (
	"context"
	"time"
)

type Parent struct {
	Name 			string 		`json:"name"`
	Address 		string 		`json:"address"`
	Profession		string 		`json:"profession"`
	Telp			string 		`json:"telp"`
}

type Familly struct {
	Status			string 		`json:"status"`
	Order			string 		`json:"order"`
}

type Student struct {
	ID				int64		`json:"id"`
	SRN				string 		`json:"srn"`
	NSRN			string 		`json:"nsrn"`
	Name			string 		`json:"name"`
	Gender			string 		`json:"gender"`
	BirthPlace		string 		`json:"birth_place"`
	BirthDate		string		`json:"birth_date"`
	Religion		Religion	`json:"religion"`
	Address			string 		`json:"address"`
	Telp			string 		`json:"telp"`
	SchoolBefore	string 		`json:"school_before"`
	AcceptedGrade	string 		`json:"accepted_grade"`
	AcceptedDate	string 		`json:"accepted_date"`
	Familly			Familly		`json:"familly"`
	Father 			Parent		`json:"fater"`
	Mother			Parent		`json:"mother"`
	Guardian		Parent		`json:"guardian"`
	CreatedAt		time.Time 	`json:"created_at"`
	UpdatedAt		time.Time 	`json:"updated_at"`
}

type StudentRepository interface {
	Fetch(ctx context.Context, cursor int64, num int64) (res []Student, err error)
}

type StudentUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Student, nextCursor string, err error)
}