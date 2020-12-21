package domain

import (
    "time"
    "context"    
)

// User ...
type User struct {
    ID          int64       `json:"id"`
    Name        string      `json:"name"`
    Email       string      `json:"email"`
    Password    string      `json:"password"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

/**
 * DTOUserLoginRequest
 * for store user request payload
 */
type DTOUserLoginRequest struct {
	Email		string 		`json:"email" validate:"required,email"`
	Password	string		`json:"password" validate:"required"`
}

// UserUsecase represent the user's usecase
type UserUsecase interface {
    Authentication(ctx context.Context, ur DTOUserLoginRequest) (t DTOTokenResponse, err error)
    RefreshToken(ctx context.Context, ur DTOTokenResponse) (t DTOTokenResponse, err error)
}

// UserRepository represent the user's repository
type UserRepository interface {
    Fetch(ctx context.Context, num int64) (res []User, err error)
    GetByEmail(ctx context.Context, email string) (res User, err error)
}

// UserCacheRepository represent the user's caching
type UserCacheRepository interface {
    StoreAuth(ctx context.Context, u User, td *TokenDetails) (err error)
    DeleteAuth(ctx context.Context, uuid string) (err error)
}