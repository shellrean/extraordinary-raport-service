package usecase

import (
    "context"
    "time"
    "log"
    "strings"

    "golang.org/x/crypto/bcrypt"

    "github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/config"
    "github.com/shellrean/extraordinary-raport/entities/helper"
)

type usecase struct {
    userRepo        domain.UserRepository
    userCacheRepo   domain.UserCacheRepository
    contextTimeout  time.Duration
    cfg             *config.Config
}

func New(d domain.UserRepository, dc domain.UserCacheRepository, timeout time.Duration, cfg *config.Config) domain.UserUsecase {
    return &usecase {
        userRepo:       d,
        userCacheRepo:  dc,
        contextTimeout: timeout,
        cfg:            cfg,
    }
}

func (u *usecase) getError(err error) (error) {
    if u.cfg.Release {
        log.Println(err.Error())
        return domain.ErrServerError
    }
    return err
}

func (u *usecase) Fetch(c context.Context, query string, cursor string, num int64) (res []domain.User, nextCursor string, err error) {
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

    res, err = u.userRepo.Fetch(ctx, query, decodedCursor, num)
    if err != nil {
        return res, nextCursor, u.getError(err)
    }

    if len(res) == int(num) {
        nextCursor = helper.EncodeCursor(res[len(res)-1].ID)
    }

    return
}

func (u *usecase) Authentication(c context.Context, ur domain.User) (td domain.Token, err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    user, err := u.userRepo.GetByEmail(ctx, ur.Email)
    if err != nil {
        if u.cfg.Release {
            log.Println(err.Error())
            return domain.Token{}, domain.ErrServerError
        }
        return domain.Token{}, err
    }

    if user == (domain.User{}) {
        return domain.Token{}, domain.ErrCredential
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ur.Password))
    if err != nil {
        return domain.Token{}, domain.ErrCredential
    }

    helper.GenerateTokenDetail(&td)

    err = helper.CreateAccessToken(u.cfg.JWTAccessKey, user, &td)
    if err != nil {
        return domain.Token{}, err
    }
    
    err = helper.CreateRefreshToken(u.cfg.JWTRefreshKey, user, &td)
    if err != nil {
        return domain.Token{}, err
    }

    if u.cfg.Redis.Enable {
        err = u.userCacheRepo.StoreAuth(ctx, user, &td)
        if err != nil {
            return domain.Token{}, err
        }
    }

    return
}

func (u *usecase) RefreshToken(c context.Context, td *domain.Token) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()
    
    token, errs := helper.VerifyToken(u.cfg.JWTRefreshKey, td.RefreshToken)
    if errs != nil {
        return domain.ErrSessVerifation
    }
    err = helper.TokenValid(token)
    if err != nil {
        return domain.ErrSessInvalid
    }

    data := helper.ExtractTokenMetadata(token)
    if u.cfg.Redis.Enable {
        err = u.userCacheRepo.DeleteAuth(ctx, data["refresh_uuid"].(string))
        if err != nil {
            return err
        }
    }
    
    user, err := u.userRepo.GetByID(ctx, int64(data["user_id"].(float64)))
    if err != nil {
        return u.getError(err)
    }
    if user == (domain.User{}) {
        return domain.ErrUserDataNotFound
    }

    helper.GenerateTokenDetail(td)

    err = helper.CreateAccessToken(u.cfg.JWTAccessKey, user, td)
    if err != nil {
        return err
    }
    
    err = helper.CreateRefreshToken(u.cfg.JWTRefreshKey, user, td)
    if err != nil {
        return err
    }

    if u.cfg.Redis.Enable {
        err = u.userCacheRepo.StoreAuth(ctx, user, td)
        if err != nil {
            return err
        }
    }

    return
}

func (u *usecase) GetByID(c context.Context, id int64) (res domain.User, err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    res, err = u.userRepo.GetByID(ctx, id)
    if err != nil {
        return res, u.getError(err)
    }
    if res == (domain.User{}) {
        return domain.User{}, domain.ErrUserDataNotFound
    }
    return
}

func (u *usecase) UserCurrentLogin(c context.Context, userID int64) (res domain.User, err error) {
    res, err = u.GetByID(c, userID)
    return
}

func (u *usecase) Store(c context.Context, ur *domain.User) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    row, err := u.userRepo.GetByEmail(ctx, ur.Email)
    if err != nil {
        return u.getError(err)
    }
    if row != (domain.User{}) {
        err = domain.ErrEmailExist
        return
    }

    password, err := bcrypt.GenerateFromPassword([]byte(ur.Password), 10)
    if err != nil {
        return u.getError(err)
    }
    ur.Password = string(password)
    ur.CreatedAt = time.Now()
    ur.UpdatedAt = time.Now()
    
    if ur.Role == 0 {
        ur.Role = domain.RoleTeacher
    }

    if err = u.userRepo.Store(ctx, ur); err != nil {
        return u.getError(err)
    }

    return
}

func (u *usecase) ImportFromExcel(c context.Context, file string) (err error) {
    ctx, cancel := context.WithTimeout(c, time.Duration(4) * time.Second)
    defer cancel()

    res, err := helper.ReadUserFileExcel(ctx, file)
    if err != nil {
        return u.getError(err)
    }

    if res == nil {
        err = domain.ErrFileNotAllowed
        return
    }

    err = u.userRepo.StoreMultiple(ctx, res)
    if err != nil {
        return u.getError(err)
    }

    return
}

func (u *usecase) Update(c context.Context, ur *domain.User) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    res, err := u.userRepo.GetByID(ctx, ur.ID)
    if err != nil {
        return u.getError(err)
    }
    if res == (domain.User{}) {
        return domain.ErrUserDataNotFound
    }

    row, err := u.userRepo.GetByEmail(ctx, ur.Email)
    if err != nil {
        return u.getError(err)
    }
    if row != (domain.User{}) && row.ID != ur.ID {
        err = domain.ErrEmailExist
        return
    }

    ur.UpdatedAt = time.Now()

    us := domain.User{
        ID:     ur.ID,
        Name:   ur.Name,
        Email:  ur.Email,
        UpdatedAt: ur.UpdatedAt,
    }

    err = u.userRepo.Update(ctx, &us)
    if err != nil {
        return u.getError(err)
    }

    return
}

func (u *usecase) Delete(c context.Context, id int64) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()
    
    res, err := u.userRepo.GetByID(ctx, id)
    if err != nil {
        return u.getError(err)
    }
    if res == (domain.User{}) {
        return domain.ErrUserDataNotFound
    }

    err = u.userRepo.Delete(ctx, id)
    if err != nil {
        return u.getError(err)
    }
    return
}

func (u *usecase) DeleteMultiple(c context.Context, query string) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    idV := strings.TrimRight(query, ",")
	idV = strings.TrimLeft(idV, ",") 
    ids := strings.Split(idV, ",")
    
    err = u.userRepo.DeleteMultiple(ctx, ids)
    if err != nil {
        return u.getError(err)
    }
    return
}