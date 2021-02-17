package usecase

import (
	"context"
	"time"
	"log"

	"golang.org/x/sync/errgroup"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
	"github.com/shellrean/extraordinary-raport/entities/helper"
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

func (u sprUsecase) getPlanByID(ctx context.Context, id int64) (res domain.ClassroomSubjectPlan, err error) {
	res, err = u.plRepo.GetByID(ctx, id)
	if err != nil {
		return res, u.getError(err)
	}

	if res == (domain.ClassroomSubjectPlan{}) {
		return res, domain.ErrSubjectPlanNotFound
	}

	return
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

	spre, err := u.sprRepo.GetByPlanIndexStudent(ctx, spr.Plan.ID, spr.Index, spr.Student.ID)
	if err != nil {
		return u.getError(err)
	}

	if spre == (domain.ClassroomSubjectPlanResult{}) {
		spr.CreatedAt = time.Now()
		spr.UpdatedAt = time.Now()
	
		err = u.sprRepo.Store(ctx, spr)
	} else {
		spr.ID = spre.ID
		spr.UpdatedAt = time.Now()

		err = u.sprRepo.Update(ctx, spr)
	}

	if err != nil {
		return u.getError(err)
	}
	
	return
}

func (u sprUsecase) FetchByPlan(c context.Context, planID int64) (res []domain.ClassroomSubjectPlanResult, err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	_, err = u.getPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	res, err = u.sprRepo.FetchByPlan(ctx, planID)
	if err != nil {
		return nil, u.getError(err)
	}
	
	return
}

func (u sprUsecase) ExportByPlan(c context.Context, planID int64) (token string, err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	plan, err := u.getPlanByID(ctx, planID)
	if err != nil {
		return "", err
	}

	res, err := u.sprRepo.FetchByPlan(ctx, planID)
	if err != nil {
		return "", u.getError(err)
	}

	students, err := u.sRepo.GetByClassroomAcademic(ctx, plan.Classroom.ID)
	if err != nil {
		return "", u.getError(err)
	}

	g, _ := errgroup.WithContext(c)
	chanResult := make(chan map[string]interface{})
	for _, row := range students {
		row := row
		plan := plan
		res := res
		g.Go(func() error {
			data := make([]uint, plan.CountPlan)
			for _, item := range res {
				if item.Student.ID == row.Student.ID {
					if item.Index < int(plan.CountPlan) {
						data[item.Index] = item.Number
					}
				}
			}

			result := map[string]interface{}{
				"nis": row.Student.SRN,
				"nama": row.Student.Name,
				"nilai": data,
				"rata": 0,
			}
			chanResult <- result
			return nil
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			return
		}
		close(chanResult)
	}()

	var datas []map[string]interface{}
	for result := range chanResult {
		datas = append(datas, result)
	}

	if err := g.Wait(); err != nil {
		return "", err
	}

	path, err := helper.WritePlanResultFileExcel(ctx, plan, datas)
	if err != nil {
		return "", u.getError(err)
	}
	
	token, err = helper.CreateFileAccessToken(u.cfg.JWTFileKey, path)
	if err != nil {
		return "", u.getError(err)
	}

	return token, nil
}