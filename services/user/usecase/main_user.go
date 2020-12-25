package usecase

import (
    "context"
    "time"

    "golang.org/x/crypto/bcrypt"

    "github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/config"
    "github.com/shellrean/extraordinary-raport/entities/helper"
)

type userUsecase struct {
    userRepo        domain.UserRepository
    userCacheRepo   domain.UserCacheRepository
    contextTimeout  time.Duration
    cfg             *config.Config
}

func NewUserUsecase(d domain.UserRepository, dc domain.UserCacheRepository, timeout time.Duration, cfg *config.Config) domain.UserUsecase {
    return &userUsecase {
        userRepo:       d,
        userCacheRepo:  dc,
        contextTimeout: timeout,
        cfg:            cfg,
    }
}

func (u *userUsecase) Fetch(c context.Context, cursor string, num int64) (result []domain.DTOUserShow, nextCursor string, err error) {
    if num == 0 {
        num = int64(10)
    }

    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    decodedCursor, err := helper.DecodeCursor(cursor)
    if err != nil && cursor != "" {
        err = domain.ErrBadParamInput
        return
    }

    var res []domain.User
    res, err = u.userRepo.Fetch(ctx, decodedCursor, num)
    if err != nil {
        return
    }

    for _, item := range res {
        user := domain.DTOUserShow{
            ID:     item.ID,
            Name:   item.Name,
            Email:  item.Email,
            CreatedAt: item.CreatedAt,
            UpdatedAt: item.UpdatedAt,
        }
        result = append(result, user)
    }

    if len(result) == int(num) {
        nextCursor = helper.EncodeCursor(result[len(result)-1].ID)
    }

    return
}

func (u *userUsecase) Authentication(c context.Context, ur domain.DTOUserLoginRequest) (t domain.DTOTokenResponse, err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    user, err := u.userRepo.GetByEmail(ctx, ur.Email)
    if err != nil {
        return
    }

    if user == (domain.User{}) {
        return domain.DTOTokenResponse{}, domain.ErrUserNotFound
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ur.Password))
    if err != nil {
        return domain.DTOTokenResponse{}, domain.ErrCredential
    }

    td := &domain.TokenDetails{}
    helper.GenerateTokenDetail(td)

    err = helper.CreateAccessToken(u.cfg.JWTAccessKey, user, td)
    if err != nil {
        return domain.DTOTokenResponse{}, err
    }
    
    err = helper.CreateRefreshToken(u.cfg.JWTRefreshKey, user, td)
    if err != nil {
        return domain.DTOTokenResponse{}, err
    }

    if u.cfg.Redis.Enable {
        err = u.userCacheRepo.StoreAuth(ctx, user, td)
        if err != nil {
            return domain.DTOTokenResponse{}, err
        }
    }

    t = domain.DTOTokenResponse{
        AccessToken:    td.AccessToken,
        RefreshToken:   td.RefreshToken,
    }

    return
}

func (u *userUsecase) RefreshToken(c context.Context, to domain.DTOTokenResponse) (t domain.DTOTokenResponse, err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()
    
    token, errs := helper.VerifyToken(u.cfg.JWTRefreshKey, to.RefreshToken)
    if errs != nil {
        return domain.DTOTokenResponse{}, domain.ErrSessVerifation
    }
    err = helper.TokenValid(token)
    if err != nil {
        return domain.DTOTokenResponse{}, domain.ErrSessInvalid
    }

    data := helper.ExtractTokenMetadata(token)
    err = u.userCacheRepo.DeleteAuth(ctx, data["refresh_uuid"].(string))
    if err != nil {
        return domain.DTOTokenResponse{}, err
    }

    user := domain.User{
        ID: int64(data["user_id"].(float64)),
    }

    td := &domain.TokenDetails{}
    helper.GenerateTokenDetail(td)

    err = helper.CreateAccessToken(u.cfg.JWTAccessKey, user, td)
    if err != nil {
        return domain.DTOTokenResponse{}, err
    }
    
    err = helper.CreateRefreshToken(u.cfg.JWTRefreshKey, user, td)
    if err != nil {
        return domain.DTOTokenResponse{}, err
    }

    if u.cfg.Redis.Enable {
        err = u.userCacheRepo.StoreAuth(ctx, user, td)
        if err != nil {
            return domain.DTOTokenResponse{}, err
        }
    }

    t = domain.DTOTokenResponse{
        AccessToken:    td.AccessToken,
        RefreshToken:   td.RefreshToken,
    }

    return
}