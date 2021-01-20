package usecase

import (
	"context"
	"time"
	"fmt"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type academicUsecase struct {
	academicRepo		domain.AcademicRepository
	contextTimeout 		time.Duration
	cfg 				*config.Config
}

func NewAcademicUsecase(d domain.AcademicRepository, timeout time.Duration, cfg *config.Config) domain.AcademicUsecase {
	return &academicUsecase {
		academicRepo:		d,
		contextTimeout:		timeout,
		cfg:				cfg,
	}	
}

func (u *academicUsecase) Fetch(c context.Context) (res []domain.Academic, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.academicRepo.Fetch(ctx)
	if err != nil {
		if u.cfg.Release {
            err = domain.ErrServerError
            return
        }
        return
	}

	return
}

func (u *academicUsecase) Generate(c context.Context) (res domain.Academic, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	currentTime := time.Now() 
	month := currentTime.Month()

	var semester uint8
	var year string
	if int(month) >= 6 {
		year = fmt.Sprintf("%d/%d", currentTime.Year(), currentTime.AddDate(1,0,0).Year())
		semester = 1
	} else {
		year = fmt.Sprintf("%d/%d", currentTime.AddDate(-1,0,0).Year(), currentTime.Year())
		semester = 2
	}
	list, err := u.academicRepo.GetByYearAndSemester(ctx, year, int(semester))
	if err != nil {
		if u.cfg.Release {
            err = domain.ErrServerError
            return
        }
        return
	}

	if list != (domain.Academic{}) {
		err = domain.ErrExistData
		return
	}

	res = domain.Academic{
		Name:	year,
		Semester: semester,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.academicRepo.Store(ctx, &res)
	if err != nil {
		if u.cfg.Release {
            err = domain.ErrServerError
            return
        }
        return
	}
	return
}