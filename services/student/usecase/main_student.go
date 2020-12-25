package usecase

import (
	"context"
	"time"

	"github.com/shellrean/extraordinary-raport/domain"
)

type studentUsecase struct {
	studentRepo		domain.StudentRepository
	contextTimeout 	time.Duration
}

func NewStudentUsecase(d domain.StudentRepository, timeout time.Duration) domain.StudentUsecase {
	return &studentUsecase {
		studentRepo:		d,
		contextTimeout:		timeout,
	}
}

func (u *studentUsecase) Fetch(c context.Context, num int64) (res []domain.Student, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.studentRepo.Fetch(ctx, num)
	if err != nil {
		return nil, err
	}

	return
}