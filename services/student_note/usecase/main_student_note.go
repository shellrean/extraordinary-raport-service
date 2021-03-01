package usecase

import (
	"context"
	"time"
	"log"

	"github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/config"
)

type usecase struct {
	snRepo 		domain.StudentNoteRepository
	csRepo 		domain.ClassroomStudentRepository
	uRepo 		domain.UserRepository
	timeout 	time.Duration
	cfg 		*config.Config
}

func New(
	m1 domain.StudentNoteRepository,
	m2 domain.ClassroomStudentRepository,
	m3 domain.UserRepository,
	timeout time.Duration, 
	cfg *config.Config,
) domain.StudentNoteUsecase {
	return &usecase {
		snRepo: 	m1,
		csRepo: 	m2,
		uRepo:		m3,
		timeout: 	timeout,
		cfg:		cfg,
	}	
}

func (u *usecase) getError(err error) (error) {
    if u.cfg.Release {
        log.Println(err.Error())
        return domain.ErrServerError
    }
    return err
}

func (u *usecase) Store(c context.Context, sn *domain.StudentNote) (err error) {
	ctx, cancel := context.WithTimeout(c, u.timeout)
	defer cancel()

	cs, err := u.csRepo.GetByID(ctx, sn.Student.ID)
	if err != nil {
		return u.getError(err)
	}

	if cs == (domain.ClassroomStudent{}) {
		return domain.ErrClassroomStudentNotFound
	}

	ue, err := u.uRepo.GetByID(ctx, sn.Teacher.ID)
	if err != nil {
		return u.getError(err)
	}

	if ue == (domain.User{}) {
		return domain.ErrUserDataNotFound
	}

	noteTypes := []int{
		domain.NoteDailly,
		domain.NoteGrd,
		domain.NoteAcademic,
		domain.NoteCharIntegrity,
		domain.NoteCharReligius,
		domain.NoteCharNation,
		domain.NoteCharIndependence,
		domain.NoteCharTeamwork,
	}

	var noteExist bool
	for _, item := range noteTypes {
		if item == sn.Type {
			noteExist = true
			break;
		}
	}

	if !noteExist {
		return domain.ErrNoteTypeNotFound
	}

	sn.CreatedAt = time.Now()
	sn.UpdatedAt = time.Now()

	err = u.snRepo.Store(ctx, sn)
	if err != nil {
		return u.getError(err)
	}
	return
}