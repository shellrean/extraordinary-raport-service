package usecase

import (
	"context"
	"time"
	"log"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type csuUsecase struct {
	csuRepo			domain.ClassroomSubjectRepository
	csaRepo 		domain.ClassroomAcademicRepository
	subjectRepo 	domain.SubjectRepository
	userRepo		domain.UserRepository
	contextTimeout  time.Duration
	cfg 			*config.Config
}

func NewClassroomSubjectUsecase(
	m 		domain.ClassroomSubjectRepository, 
	m2 		domain.ClassroomAcademicRepository, 
	m3		domain.SubjectRepository,
	m4 		domain.UserRepository,
	timeout time.Duration, 
	cfg 	*config.Config,
) domain.ClassroomSubjectUsecase{
	return &csuUsecase{
		csuRepo:		m,
		csaRepo:		m2,
		subjectRepo:	m3,
		userRepo:		m4,
		contextTimeout: timeout,
		cfg:			cfg,
	}
}

func (u *csuUsecase) FetchByClassroom(c context.Context, academicClassroomID int64) (res []domain.ClassroomSubject, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	ac, err := u.csaRepo.GetByID(ctx, academicClassroomID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
	}
	
	if ac == (domain.ClassroomAcademic{}) {
		err = domain.ErrClassroomAcademicNotFound
		return
	}

	res, err = u.csuRepo.FetchByClassroom(ctx, academicClassroomID)
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

func (u *csuUsecase) Store(c context.Context, cs *domain.ClassroomSubject) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	csa, err := u.csaRepo.GetByID(ctx, cs.ClassroomAcademic.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
	}

	if csa == (domain.ClassroomAcademic{}) {
		err = domain.ErrClassroomAcademicNotFound
		return
	}

	su, err := u.subjectRepo.GetByID(ctx, cs.Subject.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
	}

	if su == (domain.Subject{}) {
		err = domain.ErrSubjectNotFound
		return
	}

	usr, err := u.userRepo.GetByID(ctx, cs.Teacher.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
	}
	
	if usr == (domain.User{}) {
		err = domain.ErrUserDataNotFound
		return
	}

	exist, err := u.csuRepo.GetByClassroomAndSubject(ctx, cs.ClassroomAcademic.ID, cs.Subject.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
	}

	if exist != (domain.ClassroomSubject{}) {
		err = domain.ErrExistData
		return
	}

	cs.CreatedAt = time.Now()
	cs.UpdatedAt = time.Now()
	err = u.csuRepo.Store(ctx, cs)
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

func (u *csuUsecase) GetByID(c context.Context, id int64) (res domain.ClassroomSubject, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.csuRepo.GetByID(ctx, id)
	if err != nil {
		if u.cfg.Release {
			log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
	}

	if res == (domain.ClassroomSubject{}) {
		err = domain.ErrClassroomSubjectNotFound
		return
	}

	return
}

func (u *csuUsecase) Update(c context.Context, cs *domain.ClassroomSubject) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.csuRepo.GetByID(ctx, cs.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
	}

	if res == (domain.ClassroomSubject{}) {
		err = domain.ErrClassroomSubjectNotFound
		return
	}

	csa, err := u.csaRepo.GetByID(ctx, cs.ClassroomAcademic.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
	}

	if csa == (domain.ClassroomAcademic{}) {
		err = domain.ErrClassroomAcademicNotFound
		return
	}

	su, err := u.subjectRepo.GetByID(ctx, cs.Subject.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
	}

	if su == (domain.Subject{}) {
		err = domain.ErrSubjectNotFound
		return
	}

	usr, err := u.userRepo.GetByID(ctx, cs.Teacher.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
	}
	
	if usr == (domain.User{}) {
		err = domain.ErrUserDataNotFound
		return
	}

	exist, err := u.csuRepo.GetByClassroomAndSubject(ctx, cs.ClassroomAcademic.ID, cs.Subject.ID)
	if err != nil {
		if u.cfg.Release {
			log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
	}

	if exist != (domain.ClassroomSubject{}) && exist.ID != cs.ID {
		err = domain.ErrExistData
		return
	}

	cs.UpdatedAt = time.Now()
	err = u.csuRepo.Update(ctx, cs)
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

func (u *csuUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.csuRepo.GetByID(ctx, id)
	if err != nil {
		if u.cfg.Release {
			log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
	}

	if res == (domain.ClassroomSubject{}) {
		err = domain.ErrClassroomSubjectNotFound
		return
	}

	err = u.csuRepo.Delete(ctx, id)
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