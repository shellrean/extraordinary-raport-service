package usecase

import (
	"context"
	"time"
	"log"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type usecase struct {
	majorRepo		domain.MajorRepository
	contextTimeout	time.Duration
	cfg				*config.Config
}

func New(m domain.MajorRepository, timeout time.Duration, cfg *config.Config) domain.MajorUsecase {
	return &usecase {
		majorRepo:			m,
		contextTimeout: 	timeout,
		cfg:				cfg,
	}
}

func (u *usecase) getError(err error) (error) {
	if u.cfg.Release {
		log.Println(err.Error())
		return domain.ErrServerError
	}
	return err
}

func (u *usecase) Fetch(c context.Context) (res []domain.Major, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.majorRepo.Fetch(ctx)
	if err != nil {
		return res, u.getError(err)
	}
	return
}