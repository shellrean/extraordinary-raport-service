package usecase

import (
	"context"
	"time"
	"strconv"
	"log"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
)

type classroomAcademicUsecase struct {
	classAcademicRepo		domain.ClassroomAcademicRepository
	settingRepo				domain.SettingRepository
	userRepo				domain.UserRepository
	classRepo				domain.ClassroomRepository
	contextTimeout 			time.Duration
	cfg 					*config.Config
}

func NewClassroomAcademicUsecase(
	d 		domain.ClassroomAcademicRepository, 
	s 		domain.SettingRepository,
	u 		domain.UserRepository,
	c 		domain.ClassroomRepository,
	timeout time.Duration, 
	cfg 	*config.Config,
) domain.ClassroomAcademicUsecase {
	return &classroomAcademicUsecase {
		classAcademicRepo:		d,
		settingRepo:			s,
		userRepo:				u,
		classRepo:				c,
		contextTimeout:			timeout,
		cfg:					cfg,
	}
}

func (u classroomAcademicUsecase) getError(err error) (error) {
	if u.cfg.Release {
		log.Println(err.Error())
		return domain.ErrServerError
	}
	return err
}

func (u classroomAcademicUsecase) Fetch(c context.Context) (res []domain.ClassroomAcademic, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	sett, err := u.settingRepo.GetByName(ctx, domain.SettingAcademicActive)
	if err != nil {
		return nil, u.getError(err)
	}

	if sett == (domain.Setting{}) {
		err = domain.ErrNotFound
		return
	}

	id, err := strconv.Atoi(sett.Value)
	if err != nil {
		return nil, u.getError(err)
	}

	res, err = u.classAcademicRepo.Fetch(ctx, int64(id))
	if err != nil {
		return nil, u.getError(err)
	}

	return
}

func (u classroomAcademicUsecase) FetchByAcademic(c context.Context, academicID int64) (res []domain.ClassroomAcademic, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.classAcademicRepo.Fetch(ctx, academicID)
	if err != nil {
		return nil, u.getError(err)
	}

	return
}

func (u classroomAcademicUsecase) GetByID(c context.Context, id int64) (res domain.ClassroomAcademic, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	
	res, err = u.classAcademicRepo.GetByID(ctx, id)
	if err != nil {
		return domain.ClassroomAcademic{}, u.getError(err)
	}

	if res == (domain.ClassroomAcademic{}) {
		err = domain.ErrClassroomAcademicNotFound
		return
	}
	return
}

func (u classroomAcademicUsecase) Store(c context.Context, ca *domain.ClassroomAcademic) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.settingRepo.GetByName(ctx, domain.SettingAcademicActive)
	if err != nil {
		return u.getError(err)
	}

	if res == (domain.Setting{}) {
		return domain.ErrNotFound
	}

	id, err := strconv.Atoi(res.Value)
	if err != nil {
		return u.getError(err)
	}

	exist, err := u.classAcademicRepo.GetByAcademicAndClass(ctx, int64(id), ca.Classroom.ID)
	if err != nil {
		return u.getError(err)
	}

	if exist != (domain.ClassroomAcademic{}) {
		return domain.ErrClassroomAcademicExist
	}

	user, err := u.userRepo.GetByID(ctx, ca.Teacher.ID)
	if err != nil {
		return u.getError(err)
	}

	if user == (domain.User{}) {
		return domain.ErrUserDataNotFound
	}

	class, err := u.classRepo.GetByID(ctx, ca.Classroom.ID)
	if err != nil {
		return u.getError(err)
	}

	if class == (domain.Classroom{}) {
		return domain.ErrClassroomNotFound
	}

	academic := domain.Academic{
		ID: 	int64(id),
	}

	ca.Academic = academic
	ca.CreatedAt = time.Now()
	ca.UpdatedAt = time.Now()
	err = u.classAcademicRepo.Store(ctx, ca)
	if err != nil {
		return u.getError(err)
	}
	return nil
}

func (u classroomAcademicUsecase) Update(c context.Context, ca *domain.ClassroomAcademic) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	csr, err := u.classAcademicRepo.GetByID(ctx, ca.ID)
	if err != nil {
		return u.getError(err)
	}

	if csr == (domain.ClassroomAcademic{}) {
		return domain.ErrClassroomAcademicNotFound
	}

	res, err := u.settingRepo.GetByName(ctx, domain.SettingAcademicActive)
	if err != nil {
		return u.getError(err)
	}

	if res == (domain.Setting{}) {
		return domain.ErrSettingNotFound
	}

	id, err := strconv.Atoi(res.Value)
	if err != nil {
		return u.getError(err)
	}

	exist, err := u.classAcademicRepo.GetByAcademicAndClass(ctx, int64(id), ca.Classroom.ID)
	if err != nil {
		return u.getError(err)
	}

	if exist != (domain.ClassroomAcademic{}) && exist.ID != ca.ID {
		return domain.ErrClassroomAcademicExist
	}

	user, err := u.userRepo.GetByID(ctx, ca.Teacher.ID)
	if err != nil {
		return u.getError(err)
	}

	if user == (domain.User{}) {
		return domain.ErrUserDataNotFound
	}

	class, err := u.classRepo.GetByID(ctx, ca.Classroom.ID)
	if err != nil {
		return u.getError(err)
	}

	if class == (domain.Classroom{}) {
		return domain.ErrClassroomNotFound
	}

	ca.UpdatedAt = time.Now()
	err = u.classAcademicRepo.Update(ctx, ca)
	if err != nil {
		return u.getError(err)
	}
	return nil
}

func (u classroomAcademicUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.classAcademicRepo.GetByID(ctx, id)
	if err != nil {
		return u.getError(err)
	}

	if res == (domain.ClassroomAcademic{}) {
		err = domain.ErrClassroomAcademicNotFound
		return
	}

	err = u.classAcademicRepo.Delete(ctx, id)
	if err != nil {
		return u.getError(err)
	}
	return
}