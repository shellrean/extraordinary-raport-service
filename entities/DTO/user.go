package dto

type UserLogin struct {
	Email		string 		`json:"email" validate:"required,email"`
	Password	string		`json:"password" validate:"required"`
}

type UserResponse struct {
	ID 			int64 		`json:"id"`
	Name 		string 		`json:"name"`
	Email       string      `json:"email" validate:"required,email"`
}