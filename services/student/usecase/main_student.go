package usecase

import (
	"context"
	"time"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
	"github.com/shellrean/extraordinary-raport/entities/helper"
)

type studentUsecase struct {
	studentRepo		domain.StudentRepository
	contextTimeout 	time.Duration
	cfg 			*config.Config
}

func NewStudentUsecase(d domain.StudentRepository, timeout time.Duration, cfg *config.Config) domain.StudentUsecase {
	return &studentUsecase {
		studentRepo:		d,
		contextTimeout:		timeout,
		cfg: 				cfg,
	}
}

func (u *studentUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Student, nextCursor string, err error) {
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
	

	res, err = u.studentRepo.Fetch(ctx, decodedCursor, num)
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

func (u *studentUsecase) GetByID(c context.Context, id int64) (res domain.Student, err error) {
	ctx, cancel := context.WithTimeout(c,u.contextTimeout)
	defer cancel()

	res, err = u.studentRepo.GetByID(ctx, id)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}
	return
}

func (u *studentUsecase) Store(c context.Context, s *domain.Student) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()


	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	
	err = u.studentRepo.Store(ctx, s)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}

	return
}

func (u *studentUsecase) Update(c context.Context, s *domain.Student) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	s.UpdatedAt = time.Now()
	err = u.studentRepo.Update(ctx, s)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}

	return
}