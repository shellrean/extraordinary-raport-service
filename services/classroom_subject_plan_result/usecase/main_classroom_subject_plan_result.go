package usecase

import (
	"context"
	"time"
	"log"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type sprRepoT domain.ClassroomSubjectPlanResultRepository
type sRepoT domain.ClassroomStudentRepository
type suRepoT domain.ClassroomSubjectRepository
type plRepoT domain.ClassroomSubjectPlanRepository

type sprUsecase struct {
	sprRepo 	sprRepoT
	sRepo  		sRepoT
	suRepo 		suRepoT
	plRepo 		plRepoT
	timeout 	time.Duration
	cfg 		*config.Config
}

func NewClassroomSubjectPlanResultUsecase(m sprRepoT, m2 sRepoT, m3 suRepoT, m4 plRepoT, timeout time.Duration, cfg *config.Config) domain.ClassroomSubjectPlanResultUsecase {
	return &sprUsecase{
		sprRepo:	m,
		sRepo:		m2,
		suRepo:		m3,
		plRepo:		m4,
		timeout:	timeout,
		cfg:		cfg,
	}
}

func (u sprUsecase) getError(payload error) (err error) {
	if u.cfg.Release{
		log.Println(payload)
		return domain.ErrServerError
	}
	return payload
}

func (u sprUsecase) Store(c context.Context, spr *domain.ClassroomSubjectPlanResult) (err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	s, err := u.sRepo.GetByID(ctx, spr.Student.ID)
	if err != nil {
		return u.getError(err)
	}

	if s == (domain.ClassroomStudent{}) {
		return domain.ErrClassroomStudentNotFound
	}

	su, err := u.suRepo.GetByID(ctx, spr.Subject.ID)
	if err != nil {
		return u.getError(err)
	}

	if su == (domain.ClassroomSubject{}) {
		return domain.ErrClassroomSubjectNotFound
	}

	pl, err := u.plRepo.GetByID(ctx, spr.Plan.ID)
	if err != nil {
		return u.getError(err)
	}

	if pl == (domain.ClassroomSubjectPlan{}) {
		return domain.ErrSubjectPlanNotFound
	}

	spr.CreatedAt = time.Now()
	spr.UpdatedAt = time.Now()

	err = u.sprRepo.Store(ctx, spr)
	if err != nil {
		return u.getError(err)
	}
	
	return
}