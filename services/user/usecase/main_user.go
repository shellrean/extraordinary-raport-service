package usecase

import (
    "context"
    "time"

    "github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/entities/helper"
)

type userUsecase struct {
    userRepo        domain.UserRepository
    contextTimeout  time.Duration
}

func NewUserUsecase(d domain.UserRepository, timeout time.Duration) domain.UserUsecase {
    return &userUsecase {
        userRepo:       d,
        contextTimeout: timeout,
    }
}

func (u *userUsecase) Fetch(c context.Context, num int64) (res []domain.User, err error) {
    if num == 0 {
        num = 10
    }

    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    res, err = u.userRepo.Fetch(ctx, num)
    if err != nil {
        return nil, err
    }

    return
}

func (u *userUsecase) Authentication(c context.Context, ur domain.DTOUserLoginRequest, key string) (token string, err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    user, err := u.userRepo.GetByEmail(ctx, ur.Email)
    if err != nil {
        return
    }

    if user == (domain.User{}) {
        return token, domain.ErrNotFound
    }

    token, err = helper.CreateToken(key, user)

    return
}