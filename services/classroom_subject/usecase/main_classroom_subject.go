package usecase

import (
	"context"
	"time"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type csuUsecase struct {
	csuRepo			domain.ClassroomSubjectRepository
	contextTimeout  time.Duration
	cfg 			*config.Config
}

func NewClassroomSubjectUsecase(m domain.ClassroomSubjectRepository, timeout time.Duration, cfg *config.Config) domain.ClassroomSubjectUsecase{
	return &csuUsecase{
		csuRepo:		m,
		contextTimeout: timeout,
		cfg:			cfg,
	}
}

func (u *csuUsecase) FetchByClassroom(c context.Context, academicClassroomID int64) (res []domain.ClassroomSubject, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.csuRepo.FetchByClassroom(ctx, academicClassroomID)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}

	return
}