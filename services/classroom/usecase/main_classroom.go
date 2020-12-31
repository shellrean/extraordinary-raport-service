package usecase

import (
	"context"
	"time"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type classroomUsecase struct {
	classRepo 		domain.ClassroomRepository
	contextTimeout  time.Duration
	cfg 			*config.Config
}

func NewClassroomUsecase(d domain.ClassroomRepository, timeout time.Duration, cfg *config.Config) domain.ClassroomUsecase {
	return &classroomUsecase {
		classRepo:		d,
		contextTimeout:	timeout,
		cfg:			cfg,
	}
}

func (u *classroomUsecase) Fetch(c context.Context) (res []domain.Classroom, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.classRepo.Fetch(ctx)
	if err != nil {
		if u.cfg.Release {
            err = domain.ErrServerError
            return
        }
        return
	}

	return
}