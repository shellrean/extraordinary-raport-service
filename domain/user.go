package domain

import (
    "time"
    "context"    
)

// User ...
type User struct {
    ID          int64       `json:"id"`
    Name        string      `json:"name" validate:"required"`
    Email       string      `json:"email" validate:"required,email"`
    Password    string      `json:"password" validate:"required,min=6"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

// DTOUserLoginRequest
// for store user request payload
type DTOUserLoginRequest struct {
	Email		string 		`json:"email" validate:"required,email"`
	Password	string		`json:"password" validate:"required"`
}

// DTOUserList
//for response list of user
type DTOUserShow struct {
    ID          int64       `json:"id"`
    Name        string      `json:"name"`
    Email       string      `json:"email"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

// UserUsecase represent the user's usecase
type UserUsecase interface {
    Fetch(ctx context.Context, cursor string, num int64) ([]DTOUserShow, string, error)
    GetByID(ctx context.Context, id int64) (DTOUserShow, error)
    Store(ctx context.Context, ur User) (DTOUserShow, error)
    Authentication(ctx context.Context, ur DTOUserLoginRequest) (DTOTokenResponse, error)
    RefreshToken(ctx context.Context, ur DTOTokenResponse) (DTOTokenResponse, error)
}

// UserRepository represent the user's repository
type UserRepository interface {
    Fetch(ctx context.Context, cursor int64, num int64) ([]User, error)
    GetByID(ctx context.Context, id int64) (User, error)
    GetByEmail(ctx context.Context, email string) (User, error)
    Store(ctx context.Context, u *User) (error)
    Update(ctx context.Context, u *User) (error)
    Delete(ctx context.Context, id int64) (error)
}

// UserCacheRepository represent the user's caching
type UserCacheRepository interface {
    StoreAuth(ctx context.Context, u User, td *TokenDetails) (error)
    DeleteAuth(ctx context.Context, uuid string) (error)
}