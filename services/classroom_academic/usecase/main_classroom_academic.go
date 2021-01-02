package usecase

import (
	"context"
	"time"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type classroomAcademicUsecase struct {
	classAcademicRepo		domain.ClassroomAcademicRepository
	contextTimeout 			time.Duration
	cfg 					*config.Config
}

func NewClassroomAcademicUsecase(
	d 		domain.ClassroomAcademicRepository, 
	timeout time.Duration, 
	cfg 	*config.Config,
) domain.ClassroomAcademicUsecase {
	return &classroomAcademicUsecase {
		classAcademicRepo:		d,
		contextTimeout:			timeout,
		cfg:					cfg,
	}
}

func (u *classroomAcademicUsecase) Fetch(c context.Context) (res []domain.ClassroomAcademic, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.classAcademicRepo.Fetch(ctx)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}

	return
}