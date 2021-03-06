package usecase

import (
	"context"
	"time"
	"log"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type usecase struct {
	classRepo 		domain.ClassroomRepository
	majorRepo		domain.MajorRepository
	contextTimeout  time.Duration
	cfg 			*config.Config
}

func New(d domain.ClassroomRepository, m domain.MajorRepository, timeout time.Duration, cfg *config.Config) domain.ClassroomUsecase {
	return &usecase {
		classRepo:		d,
		majorRepo:		m,
		contextTimeout:	timeout,
		cfg:			cfg,
	}
}

func (u *usecase) getError(err error) (error) {
	if u.cfg.Release {
		log.Println(err)
		return domain.ErrServerError
	}
	return err
}

func (u *usecase) Fetch(c context.Context) (res []domain.Classroom, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.classRepo.Fetch(ctx)
	if err != nil {
		return nil, u.getError(err)
	}

	return
}

func (u *usecase) GetByID(c context.Context, id int64) (res domain.Classroom, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	
	res, err = u.classRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Classroom{}, u.getError(err)
	}

	if res == (domain.Classroom{}) {
		err = domain.ErrClassroomNotFound
		return
	}

	return
}

func (u *usecase) Store(c context.Context, cl *domain.Classroom) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	var row domain.Major
	row, err = u.majorRepo.GetByID(ctx, cl.Major.ID)
	if err != nil {
		return u.getError(err)
	}
	if row == (domain.Major{}) {
		err = domain.ErrNotFound
		return
	}

	cl.CreatedAt = time.Now()
	cl.UpdatedAt = time.Now()

	err = u.classRepo.Store(ctx, cl)
	if err != nil {
		return u.getError(err)
	}
	return
}

func (u *usecase) Update(c context.Context, cl *domain.Classroom) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	class, err := u.classRepo.GetByID(ctx, cl.ID)
	if err != nil {
		return u.getError(err)
	}

	if class == (domain.Classroom{}) {
		err = domain.ErrClassroomNotFound
		return
	}

	major, err := u.majorRepo.GetByID(ctx, cl.Major.ID)
	if err != nil {
		return u.getError(err)
	}
	if major == (domain.Major{}) {
		err = domain.ErrNotFound
		return
	}

	cl.UpdatedAt = time.Now()
	err = u.classRepo.Update(ctx, cl)
	if err != nil {
		return u.getError(err)
	}
	return
}

func (u *usecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	class, err := u.classRepo.GetByID(ctx, id)
	if err != nil {
		return u.getError(err)
	}

	if class == (domain.Classroom{}) {
		err = domain.ErrClassroomNotFound
		return
	}
	
	err = u.classRepo.Delete(ctx, id)
	if err != nil {
		return u.getError(err)
	}

	return
}