package usecase

import (
	"context"
	"time"
	"log"
	"fmt"
	
	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type usecase struct {
	attendRepo 		domain.AttendanceRepository
	csRepo 			domain.ClassroomStudentRepository
	timeout 		time.Duration
	cfg 			*config.Config
}

func New(m domain.AttendanceRepository, m2 domain.ClassroomStudentRepository, timeout time.Duration, cfg *config.Config) domain.AttendanceUsecase {
	return &usecase{
		attendRepo: m,
		csRepo:		m2,
		timeout:	timeout,
		cfg:		cfg,
	}
}

func (u *usecase) getError(err error) (error) {
	if u.cfg.Release {
		log.Println(err.Error())
		return domain.ErrServerError
	}
	return err
}

func (u *usecase) Fetch(c context.Context, cid int64) (res []domain.Attendance, err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	res, err = u.attendRepo.Fetch(ctx, cid)
	if err != nil {
		return res, u.getError(err)
	}
	return
}

func (u *usecase) Store(c context.Context, a *domain.Attendance) (err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()
	
	cs, err := u.csRepo.GetByID(ctx, a.Student.ID)
	if err != nil {
		return u.getError(err)
	}

	if cs == (domain.ClassroomStudent{}) {
		return domain.ErrClassroomStudentNotFound
	}

	if (!(a.Type >= 1 && a.Type <= 3)) {
		return fmt.Errorf("type attendance wrong")
	}

	exist, err := u.attendRepo.GetByStudentAndType(ctx, a.Student.ID, a.Type)
	if err != nil {
		return u.getError(err)
	}

	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	if exist != (domain.Attendance{}) {
		err = u.attendRepo.Update(ctx, a)
	} else {
		err = u.attendRepo.Store(ctx, a)
	}

	if err != nil {
		return u.getError(err)
	}

	return
}