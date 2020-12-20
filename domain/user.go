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

// UserUsecase represent the user's usecase
type UserUsecase interface {
    
}

// UserRepository represent the user's repository
type UserRepository interface {
    Fetch(ctx context.Context, num int64) (res []User, err error)
    GetByEmail(ctx context.Context, email string) (res User, err error)
}