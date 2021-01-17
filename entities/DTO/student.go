package dto

type StudentResponse struct {
	ID 				int64 		`json:"id"`
	SRN 			string 		`json:"srn" validate:"required"`
	NSRN 			string 		`json:"nsrn" validate:"required"`
	Name 			string 		`json:"name" validate:"required"`
	Gender 			string 		`json:"gender"`
	BirthPlace 		string 		`json:"birth_place"`
	BirthDate		string		`json:"birth_date"`
	ReligionID  	int64 		`json:"religion_id"`
	Address 		string 		`json:"address"`
	Telp 			string 		`json:"telp"`
	SchoolBefore	string 		`json:"school_before"`
	AcceptedGrade	string 		`json:"accepted_grade"`
	AcceptedDate	string 		`json:"accepted_date"`
	FamillyStatus 	string 		`json:"familly_status"`
	FamillyOrder 	string 		`json:"familly_order"`
	FatherName 		string 		`json:"father_name"`
	FatherAddress 	string 		`json:"father_address"`
	FatherProfession string 	`json:"father_profession"`
	FatherTelp 		string 		`json:"father_telp"`
	MotherName 		string 		`json:"mother_name"`
	MotherAddress	string 		`json:"mother_address"`
	MotherProfession string 	`json:"mother_profession"`
	MotherTelp 		string 		`json:"mother_telp"`
	GrdName 		string 		`json:"grd_name"`
	GrdAddress 		string 		`json:"grd_address"`
	GrdProfession 	string 		`json:"grd_profession"`
	GrdTelp			string 		`json:"grd_telp"`
}