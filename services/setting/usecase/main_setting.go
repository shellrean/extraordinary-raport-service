package usecase

import (
	"context"
	"time"
	"log"

	"github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/config"
)

type settingUsecase struct {
	settRepo		domain.SettingRepository
	contextTimeout  time.Duration
    cfg             *config.Config
}

func NewSettingUsecase(d domain.SettingRepository, timeout time.Duration, cfg *config.Config) domain.SettingUsecase {
	return &settingUsecase {
		settRepo:		d,
		contextTimeout:	timeout,
		cfg:			cfg,
	}
}

func (u *settingUsecase) Fetch(c context.Context, names []string) (res []domain.Setting, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	
	res, err = u.settRepo.Fetch(ctx, names)
    if err != nil {
        if u.cfg.Release {
			log.Println(err.Error())
            err = domain.ErrServerError
            return
        }
        return
	}
	
	return
}