package usecase

import (
	"context"
	"time"
	"log"
	"strings"
	"strconv"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type usecase struct {
	cspRepo domain.ClassroomSubjectPlanRepository
	usrRepo domain.UserRepository
	csRepo 	domain.ClassroomSubjectRepository
	caRepo 	domain.ClassroomAcademicRepository
	settingRepo domain.SettingRepository
	timeout time.Duration
	cfg 	*config.Config
}

func New(m domain.ClassroomSubjectPlanRepository, m2 domain.UserRepository, m3 domain.ClassroomSubjectRepository, m4 domain.ClassroomAcademicRepository, m5 domain.SettingRepository, timeout time.Duration, cfg *config.Config) domain.ClassroomSubjectPlanUsecase{
	return &usecase{
		cspRepo:	m,
		usrRepo:	m2,
		csRepo:		m3,
		caRepo:		m4,
		settingRepo:m5,
		timeout:	timeout,
		cfg:		cfg,
	}
}

func (u *usecase) getError(err error) (error) {
	if u.cfg.Release {
		log.Println(err)
		return domain.ErrServerError
	}
	return err
}

func (u *usecase) Fetch(c context.Context, query string, userID int64, classID int64) (res []domain.ClassroomSubjectPlan, err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()
	

	if userID != 0 && classID != 0 {
		res, err = u.cspRepo.FetchByTeacherAndClassroom(ctx, userID, classID)
		if err != nil {
			return nil, u.getError(err)
		}
	} else if userID != 0 {
		sett, err := u.settingRepo.GetByName(ctx, domain.SettingAcademicActive)
		if err != nil {
			return nil, u.getError(err)
		}

		if sett == (domain.Setting{}) {
			return nil, domain.ErrNotFound
		}

		id, err := strconv.Atoi(sett.Value)
		if err != nil {
			return nil, u.getError(err)
		}
		
		res, err = u.cspRepo.FetchByAcademicTeacher(ctx, int64(id), userID)
		if err != nil {
			return nil, u.getError(err)
		}
		
	} else if classID != 0 {
		res, err = u.cspRepo.FetchByClassroom(ctx, classID)
		if err != nil {
			return nil, u.getError(err)
		}
	} else {
		return nil, domain.ErrBadParamInput
	}

	return
}

func (u *usecase) GetByID(c context.Context, id int64) (res domain.ClassroomSubjectPlan, err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	res, err = u.cspRepo.GetByID(ctx, id)
	if err != nil {
		return domain.ClassroomSubjectPlan{}, u.getError(err)
	}

	if res == (domain.ClassroomSubjectPlan{}) {
		return domain.ClassroomSubjectPlan{}, domain.ErrSubjectPlanNotFound
	}

	return
}

func (u *usecase) Store(c context.Context, csp *domain.ClassroomSubjectPlan) (err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	if csp.Type != domain.PlanTask && csp.Type != domain.PlanEVMS && csp.Type != domain.PlanEVLS && csp.Type != domain.PlanExam{
		err = domain.ErrValidation
		return
	}

	usr, err := u.usrRepo.GetByID(ctx, csp.Teacher.ID)
	if err != nil {
		return u.getError(err)
	}

	if usr == (domain.User{}) {
		err = domain.ErrUserDataNotFound
		return
	}

	cs, err := u.csRepo.GetByID(ctx, csp.Subject.ID)
	if err != nil {
		return u.getError(err)
	}
	
	if cs == (domain.ClassroomSubject{}) {
		err = domain.ErrClassroomSubjectNotFound
		return
	}

	ca, err := u.caRepo.GetByID(ctx, csp.Classroom.ID)
	if err != nil {
		return u.getError(err)
	}

	if ca == (domain.ClassroomAcademic{}) {
		err = domain.ErrClassroomAcademicNotFound
		return
	}

	csp.CreatedAt = time.Now()
	csp.UpdatedAt = time.Now()

	err = u.cspRepo.Store(ctx, csp)
	if err != nil {
		return u.getError(err)
	}
	
	return
}

func (u *usecase) Update(c context.Context, csp *domain.ClassroomSubjectPlan) (err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	row, err := u.cspRepo.GetByID(ctx, csp.ID)
	if err != nil {
		return u.getError(err)
	}

	if row == (domain.ClassroomSubjectPlan{}) {
		err = domain.ErrSubjectPlanNotFound
		return
	}

	if csp.Type != domain.PlanTask && csp.Type != domain.PlanEVMS && csp.Type != domain.PlanEVLS && csp.Type != domain.PlanExam{
		err = domain.ErrValidation
		return
	}

	usr, err := u.usrRepo.GetByID(ctx, csp.Teacher.ID)
	if err != nil {
		return u.getError(err)
	}

	if usr == (domain.User{}) {
		err = domain.ErrUserDataNotFound
		return
	}

	cs, err := u.csRepo.GetByID(ctx, csp.Subject.ID)
	if err != nil {
		return u.getError(err)
	}
	
	if cs == (domain.ClassroomSubject{}) {
		err = domain.ErrClassroomSubjectNotFound
		return
	}

	ca, err := u.caRepo.GetByID(ctx, csp.Classroom.ID)
	if err != nil {
		return u.getError(err)
	}

	if ca == (domain.ClassroomAcademic{}) {
		err = domain.ErrClassroomAcademicNotFound
		return
	}

	csp.UpdatedAt = time.Now()

	err = u.cspRepo.Update(ctx, csp)
	if err != nil {
		return u.getError(err)
	}
	
	return
}

func (u *usecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	csp, err := u.cspRepo.GetByID(ctx, id)
	if err != nil {
		return u.getError(err)
	}

	if csp == (domain.ClassroomSubjectPlan{}) {
		err = domain.ErrSubjectPlanNotFound
		return
	}

	err = u.cspRepo.Delete(ctx, id)
	if err != nil {
		return u.getError(err)
	}

	return	
}

func (u *usecase) DeleteMultiple(c context.Context, query string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	idV := strings.TrimRight(query, ",")
	idV = strings.TrimLeft(idV, ",") 
	ids := strings.Split(idV, ",")
	
	err = u.cspRepo.DeleteMultiple(ctx, ids)
    if err != nil {
        return u.getError(err)
    }
    return
}