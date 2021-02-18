package usecase

import (
	"context"
    "time"
    "log"

	"github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/config"
    "github.com/shellrean/extraordinary-raport/entities/helper"
)

type usecase struct {
	subjectRepo 	domain.SubjectRepository
	contextTimeout  time.Duration
    cfg             *config.Config
}

func New(
	d domain.SubjectRepository,
	timeout time.Duration,
	cfg *config.Config,
) domain.SubjectUsecase {
	return &usecase {
		subjectRepo:	d,
		contextTimeout: timeout,
		cfg:			cfg,
	}
}

func (u *usecase) getError(err error) (error) {
    if u.cfg.Release {
        log.Println(err.Error())
        return domain.ErrServerError
    }
    return err
}

func (u *usecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Subject, nextCursor string, err error) {
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
	
	res, err = u.subjectRepo.Fetch(ctx, decodedCursor, num)
    if err != nil {
        return res, nextCursor, u.getError(err)
    }

    if len(res) == int(num) {
        nextCursor = helper.EncodeCursor(res[len(res)-1].ID)
    }

    return
}

func (u *usecase) GetByID(c context.Context, id int64) (res domain.Subject, err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    res, err = u.subjectRepo.GetByID(ctx, id)
    if err != nil {
        return res, u.getError(err)
    }
    if res == (domain.Subject{}) {
        err= domain.ErrSubjectNotFound
        return
    }
    return
}

func (u *usecase) Store(c context.Context, s *domain.Subject) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    s.CreatedAt = time.Now()
    s.UpdatedAt = time.Now()

    err = u.subjectRepo.Store(ctx, s)
    if err != nil {
        return u.getError(err)
    }
    return
}

func (u *usecase) Update(c context.Context, s *domain.Subject) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    res, err := u.subjectRepo.GetByID(ctx, s.ID)
    if err != nil {
        return u.getError(err)
    }
    if res == (domain.Subject{}) {
        err= domain.ErrSubjectNotFound
        return
    }

    s.UpdatedAt = time.Now()

    err = u.subjectRepo.Update(ctx, s)
    if err != nil {
        return u.getError(err)
    }
    return 
}

func (u *usecase) Delete(c context.Context, id int64) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    res, err := u.subjectRepo.GetByID(ctx, id)
    if err != nil {
        return u.getError(err)
    }
    if res == (domain.Subject{}) {
        err= domain.ErrSubjectNotFound
        return
    }

    err = u.subjectRepo.Delete(ctx, id)
    if err != nil {
        return u.getError(err)
    }
    return
}