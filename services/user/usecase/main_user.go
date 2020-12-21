package usecase

import (
    "context"
    "time"

    "golang.org/x/crypto/bcrypt"
    "github.com/twinj/uuid"

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

func (u *userUsecase) Authentication(c context.Context, ur domain.DTOUserLoginRequest, key string) (t domain.DTOTokenResponse, err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    user, err := u.userRepo.GetByEmail(ctx, ur.Email)
    if err != nil {
        return
    }

    if user == (domain.User{}) {
        return domain.DTOTokenResponse{}, domain.ErrInvalidUser
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ur.Password))
    if err != nil {
        return domain.DTOTokenResponse{}, domain.ErrInvalidUser
    }

    td := &domain.TokenDetails{}
    td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
    td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
    td.AccessUuid = uuid.NewV4().String()
    td.RefreshUuid = uuid.NewV4().String()

    err = helper.CreateAccessToken(key, user, td)
    if err != nil {
        return domain.DTOTokenResponse{}, err
    }
    
    err = helper.CreateRefreshToken(key, user, td)
    if err != nil {
        return domain.DTOTokenResponse{}, err
    }

    t = domain.DTOTokenResponse{
        AccessToken:    td.AccessToken,
        RefreshToken:   td.RefreshToken,
    }

    return
}