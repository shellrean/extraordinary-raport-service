package usecase 

import (
    "time"
    "context"

    "github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/config"
    "github.com/shellrean/extraordinary-raport/entities/helper"
)

type csUsecase struct {
	csRepo			domain.ClassroomStudentRepository
	contextTimeout	time.Duration
	cfg 			*config.Config
}

func NewClassroomStudentUsecase(
	d domain.ClassroomStudentRepository,
	timeout time.Duration,
	cfg *config.Config,
) domain.ClassroomStudentUsecase {
	return &csUsecase {
		csRepo:			d,
		contextTimeout: timeout,
		cfg:			cfg,
	}
}

func (u *csUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.ClassroomStudent, nextCursor string, err error) {
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

    res, err = u.csRepo.Fetch(ctx, decodedCursor, num)
    if err != nil {
        if u.cfg.Release {
            err = domain.ErrServerError
            return
        }
        return
    }

    if len(res) == int(num) {
        nextCursor = helper.EncodeCursor(res[len(res)-1].ID)
    }

    return
}

func (u *csUsecase) GetByID(c context.Context, id int64) (res domain.ClassroomStudent, err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    res, err = u.csRepo.GetByID(ctx, id)
    if err != nil {
        if u.cfg.Release {
            err = domain.ErrServerError
            return
        }
        return 
    }
    if res == (domain.ClassroomStudent{}) {
        return domain.ClassroomStudent{}, domain.ErrNotFound
    }
    return
}

func (u *csUsecase) Store(c context.Context, cs *domain.ClassroomStudent) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    cs.CreatedAt = time.Now()
    cs.UpdatedAt = time.Now()

    if err = u.csRepo.Store(ctx, cs); err != nil {
        if u.cfg.Release {
            return domain.ErrServerError
        }
        return
    }

    return
}

func (u *csUsecase) Update(c context.Context, cs *domain.ClassroomStudent) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    cs.UpdatedAt = time.Now()

    if err = u.csRepo.Update(ctx, cs); err != nil {
        if u.cfg.Release {
            return domain.ErrServerError
        }
        return
    }

    return
}