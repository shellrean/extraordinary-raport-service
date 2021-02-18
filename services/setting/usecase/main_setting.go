package usecase

import (
	"context"
	"time"
	"log"
	"strconv"

	"github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/config"
)

type usecase struct {
	settRepo		domain.SettingRepository
	academicRepo	domain.AcademicRepository
	contextTimeout  time.Duration
    cfg             *config.Config
}

func New(d domain.SettingRepository, m domain.AcademicRepository, timeout time.Duration, cfg *config.Config) domain.SettingUsecase {
	return &usecase {
		settRepo:		d,
		academicRepo:	m,
		contextTimeout:	timeout,
		cfg:			cfg,
	}
}

func (u *usecase) getError(err error) (error) {
	if u.cfg.Release {
		log.Println(err.Error())
		return domain.ErrServerError
	}
	return err
}

func (u *usecase) Fetch(c context.Context, names []string) (res []domain.Setting, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	
	res, err = u.settRepo.Fetch(ctx, names)
    if err != nil {
        return res, u.getError(err)
	}
	
	return
}

func (u *usecase) beforeUpdateAcademicActive(ctx context.Context, s *domain.Setting) (err error) {
	id, err := strconv.Atoi(s.Value)
	if err != nil {
		err = domain.ErrSettingNotFound
		return
	}

	res, err := u.academicRepo.GetByID(ctx, int64(id))
	if err != nil {
		return u.getError(err)
	}

	if res == (domain.Academic{}) {
		err = domain.ErrAcademicNotFound
		return
	}

	return
}

func (u *usecase) Update(c context.Context, s *domain.Setting) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	switch s.Name {
	case domain.SettingAcademicActive:
		err = u.beforeUpdateAcademicActive(ctx, s)
		if err != nil {
			return
		}
	default:
		err = domain.ErrSettingNotFound
		return
	}

	s.UpdatedAt = time.Now()
	err = u.settRepo.Update(ctx, s)
    if err != nil {
        return u.getError(err)
	}
	
	return
}