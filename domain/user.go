package domain

import (
    "time"
    "context"    
)

type User struct {
    ID          int64
    Name        string
    Email       string
    Password    string
    Role        int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type UserUsecase interface {
    Fetch(ctx context.Context, query string, cursor string, num int64) ([]User, string, error)
    GetByID(ctx context.Context, id int64) (User, error)
    Store(ctx context.Context, ur *User) (error)
    Update(ctx context.Context, ur *User) (error)
    Delete(ctx context.Context, id int64) (error)
    DeleteMultiple(ctx context.Context, query string) (error)
    Authentication(ctx context.Context, ur User) (Token, error)
    RefreshToken(ctx context.Context, ur *Token) (error)
}

type UserRepository interface {
    Fetch(ctx context.Context, query string, cursor int64, num int64) ([]User, error)
    GetByID(ctx context.Context, id int64) (User, error)
    GetByEmail(ctx context.Context, email string) (User, error)
    Store(ctx context.Context, u *User) (error)
    Update(ctx context.Context, u *User) (error)
    Delete(ctx context.Context, id int64) (error)
    DeleteMultiple(ctx context.Context, ids []string) (error)
}

type UserCacheRepository interface {
    StoreAuth(ctx context.Context, u User, td *Token) (error)
    DeleteAuth(ctx context.Context, uuid string) (error)
}