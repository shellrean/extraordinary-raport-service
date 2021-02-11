package usecase

import (
	"context"
	"time"
	"log"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type csPlanUsecase struct {
	cspRepo domain.ClassroomSubjectPlanRepository
	usrRepo domain.UserRepository
	csRepo 	domain.ClassroomSubjectRepository
	caRepo 	domain.ClassroomAcademicRepository
	timeout time.Duration
	cfg 	*config.Config
}

func NewClassroomSubjectPlanUsecase(m domain.ClassroomSubjectPlanRepository, m2 domain.UserRepository, m3 domain.ClassroomSubjectRepository, m4 domain.ClassroomAcademicRepository, timeout time.Duration, cfg *config.Config) domain.ClassroomSubjectPlanUsecase{
	return &csPlanUsecase{
		cspRepo:	m,
		usrRepo:	m2,
		csRepo:		m3,
		caRepo:		m4,
		timeout:	timeout,
		cfg:		cfg,
	}
}

func (u csPlanUsecase) getError(err error) (error) {
	if u.cfg.Release {
		log.Println(err)
		return domain.ErrServerError
	}
	return err
}

func (u csPlanUsecase) Fetch(c context.Context, query string, userID int64, classID int64) (res []domain.ClassroomSubjectPlan, err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	if userID != 0 && classID != 0 {
		res, err = u.cspRepo.FetchByTeacherAndClassroom(ctx, userID, classID)
		if err != nil {
			return nil, u.getError(err)
		}
	} else if userID != 0 {
		res, err = u.cspRepo.FetchByTeacher(ctx, userID)
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

func (u csPlanUsecase) Store(c context.Context, csp *domain.ClassroomSubjectPlan) (err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	if csp.Type != domain.PlanTask && csp.Type != domain.PlanEVMS && csp.Type != domain.PlanEVLS && csp.Type != domain.PlanExam{
		err = domain.ErrValidation
		return
	}

	usr, err := u.usrRepo.GetByID(ctx, csp.Teacher.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err)
			err = domain.ErrServerError
			return
		}
		return
	}

	if usr == (domain.User{}) {
		err = domain.ErrUserDataNotFound
		return
	}

	cs, err := u.csRepo.GetByID(ctx, csp.Subject.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err)
			err = domain.ErrServerError
			return
		}
		return
	}
	
	if cs == (domain.ClassroomSubject{}) {
		err = domain.ErrClassroomSubjectNotFound
		return
	}

	ca, err := u.caRepo.GetByID(ctx, csp.Classroom.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err)
			err = domain.ErrServerError
			return
		}
		return
	}

	if ca == (domain.ClassroomAcademic{}) {
		err = domain.ErrClassroomAcademicNotFound
		return
	}

	csp.CreatedAt = time.Now()
	csp.UpdatedAt = time.Now()

	err = u.cspRepo.Store(ctx, csp)
	if err != nil {
		if u.cfg.Release {
			log.Println(err)
			err = domain.ErrServerError
			return
		}
		return
	}
	
	return
}

func (u csPlanUsecase) Update(c context.Context, csp *domain.ClassroomSubjectPlan) (err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	row, err := u.cspRepo.GetByID(ctx, csp.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err)
			err = domain.ErrServerError
			return
		}
		return
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
		if u.cfg.Release {
			log.Println(err)
			err = domain.ErrServerError
			return
		}
		return
	}

	if usr == (domain.User{}) {
		err = domain.ErrUserDataNotFound
		return
	}

	cs, err := u.csRepo.GetByID(ctx, csp.Subject.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err)
			err = domain.ErrServerError
			return
		}
		return
	}
	
	if cs == (domain.ClassroomSubject{}) {
		err = domain.ErrClassroomSubjectNotFound
		return
	}

	ca, err := u.caRepo.GetByID(ctx, csp.Classroom.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err)
			err = domain.ErrServerError
			return
		}
		return
	}

	if ca == (domain.ClassroomAcademic{}) {
		err = domain.ErrClassroomAcademicNotFound
		return
	}

	csp.UpdatedAt = time.Now()

	err = u.cspRepo.Update(ctx, csp)
	if err != nil {
		if u.cfg.Release {
			log.Println(err)
			err = domain.ErrServerError
			return
		}
		return
	}
	
	return
}

func (u csPlanUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	csp, err := u.cspRepo.GetByID(ctx, id)
	if err != nil {
		if u.cfg.Release {
			log.Println(err)
			err = domain.ErrServerError
			return
		}
		return
	}

	if csp == (domain.ClassroomSubjectPlan{}) {
		err = domain.ErrSubjectPlanNotFound
		return
	}

	err = u.cspRepo.Delete(ctx, id)
	if err != nil {
		if u.cfg.Release {
			log.Println(err)
			err = domain.ErrServerError
			return
		}
		return
	}

	return	
}