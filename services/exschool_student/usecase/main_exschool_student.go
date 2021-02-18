package usecase

import (
	"context"
	"time"
	"log"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type usecase struct {
	exsRepo 		domain.ExschoolStudentRepository
	exRepo 			domain.ExschoolRepository
	csRepo 			domain.ClassroomStudentRepository
	ctxTimeout 		time.Duration
	cfg 			*config.Config
}

func New(m domain.ExschoolStudentRepository, m2 domain.ExschoolRepository, m3 domain.ClassroomStudentRepository, timeout time.Duration, cfg *config.Config) domain.ExschoolStudentUsecase{
	return &usecase {
		exsRepo:	 m,
		exRepo:		 m2,
		csRepo:		 m3,
		ctxTimeout:  timeout,
		cfg:		 cfg,
	}
}

func (u *usecase) getError(err error) (error) {
    if u.cfg.Release {
        log.Println(err.Error())
        return domain.ErrServerError
    }
    return err
}

func (u *usecase) FetchByClassroom(c context.Context, id int64) (res []domain.ExschoolStudent, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	ex, err := u.exRepo.GetByID(ctx, id)
	if err != nil {
		return nil, u.getError(err)
	}

	if ex == (domain.Exschool{}) {
		return nil, domain.ErrExschoolNotFound
	}

	res, err = u.exsRepo.FetchByClassroom(ctx, id)
	if err != nil {
		return nil, u.getError(err)
	}

	return
}

func (u *usecase) GetByID(c context.Context, id int64) (res domain.ExschoolStudent, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	res, err = u.exsRepo.GetByID(ctx, id)
	if err != nil {
		return domain.ExschoolStudent{}, u.getError(err)
	}

	return
}

func (u *usecase) Store(c context.Context, exs *domain.ExschoolStudent) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	ex, err := u.exRepo.GetByID(ctx, exs.Exschool.ID)
	if err != nil {
		return u.getError(err)
	}

	if ex == (domain.Exschool{}) {
		return domain.ErrExschoolNotFound
	}

	cs, err := u.csRepo.GetByID(ctx, exs.Student.ID)
	if err != nil {
		return u.getError(err)
	}

	if cs == (domain.ClassroomStudent{}) {
		return domain.ErrClassroomStudentNotFound
	}

	exs.CreatedAt = time.Now()
	exs.UpdatedAt = time.Now()

	err = u.exsRepo.Store(ctx, exs)
	if err != nil {
		return u.getError(err)
	}

	return
}

func (u *usecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	exs, err := u.GetByID(c, id)
	if err != nil {
		return u.getError(err)
	}

	if exs == (domain.ExschoolStudent{}) {
		return domain.ErrExschoolStudentNotFound
	}

	err = u.exsRepo.Delete(ctx, id)
	if err != nil {
		return u.getError(err)
	}

	return
}