package usecase 

import (
	"context"
	"time"
	"log"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type usecase struct {
	exschoolRepo 		domain.ExschoolRepository
	ctxTimeout 			time.Duration
	cfg 				*config.Config
}

func New(d domain.ExschoolRepository, timeout time.Duration, cfg *config.Config) domain.ExschoolUsecase {
	return &usecase {
		exschoolRepo: 	d,
		ctxTimeout:		timeout,
		cfg: 			cfg,
	}
}

func (u *usecase) getError(err error) (error) {
	if u.cfg.Release {
		log.Println(err.Error())
		return domain.ErrServerError
	}
	return err
}

func (u *usecase) Fetch(c context.Context) (res []domain.Exschool, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	res, err = u.exschoolRepo.Fetch(ctx)
	if err != nil {
		return nil, u.getError(err)
	}

	return
}

func (u *usecase) GetByID(c context.Context, id int64) (res domain.Exschool, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	res, err = u.exschoolRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Exschool{}, u.getError(err)
	}

	if res == (domain.Exschool{}) {
		return domain.Exschool{}, domain.ErrExschoolNotFound
	}

	return
}

func (u *usecase) Store(c context.Context, ex *domain.Exschool) (err error) {
	return
}