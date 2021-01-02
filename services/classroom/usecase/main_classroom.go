package usecase

import (
	"context"
	"time"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type classroomUsecase struct {
	classRepo 		domain.ClassroomRepository
	majorRepo		domain.MajorRepository
	contextTimeout  time.Duration
	cfg 			*config.Config
}

func NewClassroomUsecase(d domain.ClassroomRepository, m domain.MajorRepository, timeout time.Duration, cfg *config.Config) domain.ClassroomUsecase {
	return &classroomUsecase {
		classRepo:		d,
		majorRepo:		m,
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

func (u *classroomUsecase) GetByID(c context.Context, id int64) (res domain.Classroom, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	
	res, err = u.classRepo.GetByID(ctx, id)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}
	return
}

func (u *classroomUsecase) Store(c context.Context, cl *domain.Classroom) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	var row domain.Major
	row, err = u.majorRepo.GetByID(ctx, cl.Major.ID)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}
	if row == (domain.Major{}) {
		err = domain.ErrNotFound
		return
	}

	cl.CreatedAt = time.Now()
	cl.UpdatedAt = time.Now()

	err = u.classRepo.Store(ctx, cl)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}
	return
}

func (u *classroomUsecase) Update(c context.Context, cl *domain.Classroom) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	var row domain.Major
	row, err = u.majorRepo.GetByID(ctx, cl.Major.ID)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}
	if row == (domain.Major{}) {
		err = domain.ErrNotFound
		return
	}

	cl.UpdatedAt = time.Now()
	err = u.classRepo.Update(ctx, cl)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}
	return
}

func (u *classroomUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	err = u.classRepo.Delete(ctx, id)
	if err != nil {
		if (u.cfg.Release) {
			err = domain.ErrServerError
			return
		}
		return
	}

	return
}