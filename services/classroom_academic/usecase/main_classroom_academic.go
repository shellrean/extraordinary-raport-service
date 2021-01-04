package usecase

import (
	"context"
	"time"
	"strconv"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type classroomAcademicUsecase struct {
	classAcademicRepo		domain.ClassroomAcademicRepository
	settingRepo				domain.SettingRepository
	userRepo				domain.UserRepository
	classRepo				domain.ClassroomRepository
	contextTimeout 			time.Duration
	cfg 					*config.Config
}

func NewClassroomAcademicUsecase(
	d 		domain.ClassroomAcademicRepository, 
	s 		domain.SettingRepository,
	u 		domain.UserRepository,
	c 		domain.ClassroomRepository,
	timeout time.Duration, 
	cfg 	*config.Config,
) domain.ClassroomAcademicUsecase {
	return &classroomAcademicUsecase {
		classAcademicRepo:		d,
		settingRepo:			s,
		userRepo:				u,
		classRepo:				c,
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

func (u *classroomAcademicUsecase) Store(c context.Context, ca *domain.ClassroomAcademic) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.settingRepo.GetByName(ctx, domain.SettingAcademicActive)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}

	if res == (domain.Setting{}) {
		err = domain.ErrNotFound
		return
	}

	user, err := u.userRepo.GetByID(ctx, ca.Teacher.ID)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}

	if user == (domain.User{}) {
		err = domain.ErrNotFound
		return
	}

	class, err := u.classRepo.GetByID(ctx, ca.Classroom.ID)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}

	if class == (domain.Classroom{}) {
		err = domain.ErrNotFound
		return
	}


	id, err := strconv.Atoi(res.Value)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}

	academic := domain.Academic{
		ID: 	int64(id),
	}

	ca.Academic = academic
	err = u.classAcademicRepo.Store(ctx, ca)
	if err != nil {
		if u.cfg.Release {
			err = domain.ErrServerError
			return
		}
		return
	}
	return
}