package usecase 

import (
    "time"
    "context"
    "strconv"
    "log"

    "github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/config"
    "github.com/shellrean/extraordinary-raport/entities/helper"
)

type usecase struct {
    csRepo			domain.ClassroomStudentRepository
    csaRepo         domain.ClassroomAcademicRepository
    settingRepo     domain.SettingRepository
	contextTimeout	time.Duration
	cfg 			*config.Config
}

func New(
    d domain.ClassroomStudentRepository,
    m domain.ClassroomAcademicRepository,
    sr domain.SettingRepository,
	timeout time.Duration,
	cfg *config.Config,
) domain.ClassroomStudentUsecase {
	return &usecase {
        csRepo:			d,
        csaRepo:        m,
        settingRepo:    sr,
		contextTimeout: timeout,
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

func (u *usecase) Fetch(c context.Context, cursor string, num int64) (res []domain.ClassroomStudent, nextCursor string, err error) {
    if num == 0 {
        num = int64(10)
    }

    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    decodedCursor, err := helper.DecodeCursor(cursor)
    if err != nil && cursor != "" {
        err = domain.ErrBadParamInput
        return
    }

    res, err = u.csRepo.Fetch(ctx, decodedCursor, num)
    if err != nil {
        return res, nextCursor, u.getError(err)
    }

    if len(res) == int(num) {
        nextCursor = helper.EncodeCursor(res[len(res)-1].ID)
    }

    return
}

func (u *usecase) GetByClassroomAcademic(c context.Context, classroomAcademicID int64) (res []domain.ClassroomStudent, err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    csa, err := u.csaRepo.GetByID(ctx, classroomAcademicID)
    if err != nil {
        return res, u.getError(err)
    }

    if csa == (domain.ClassroomAcademic{}) {
        err = domain.ErrClassroomAcademicNotFound
        return
    }

    res, err = u.csRepo.GetByClassroomAcademic(ctx, classroomAcademicID)
    if err != nil {
        return res, u.getError(err)
    }

    return
}

func (u *usecase) GetByID(c context.Context, id int64) (res domain.ClassroomStudent, err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    res, err = u.csRepo.GetByID(ctx, id)
    if err != nil {
        return res, u.getError(err)
    }
    if res == (domain.ClassroomStudent{}) {
        err = domain.ErrClassroomStudentNotFound
        return
    }
    return
}

func (u *usecase) Store(c context.Context, cs *domain.ClassroomStudent) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    res, err := u.settingRepo.GetByName(ctx, domain.SettingAcademicActive)
	if err != nil {
		return u.getError(err)
    }
    
    if res == (domain.Setting{}) {
		err = domain.ErrSettingNotFound
		return
    }
    
    id, err := strconv.Atoi(res.Value)
	if err != nil {
		return u.getError(err)
    }
    
    exist, err := u.csRepo.GetByAcademicAndStudent(ctx, int64(id), cs.Student.ID)
	if err != nil {
		return u.getError(err)
	}

	if exist != (domain.ClassroomStudent{}) {
        err = domain.ErrClassroomStudentExist
        return
	}

    cs.CreatedAt = time.Now()
    cs.UpdatedAt = time.Now()

    if err = u.csRepo.Store(ctx, cs); err != nil {
        return u.getError(err)
    }

    return
}

func (u *usecase) Update(c context.Context, cs *domain.ClassroomStudent) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    exist, err := u.csRepo.GetByID(ctx, cs.ID)
    if err != nil {
        return u.getError(err)
    }

    if exist == (domain.ClassroomStudent{}) {
        err = domain.ErrClassroomStudentNotFound
        return
    }

    cs.UpdatedAt = time.Now()

    if err = u.csRepo.Update(ctx, cs); err != nil {
        return u.getError(err)
    }

    return
}

func (u *usecase) Delete(c context.Context, id int64) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    exist, err := u.csRepo.GetByID(ctx, id)
    if err != nil {
        return u.getError(err)
    }

    if exist == (domain.ClassroomStudent{}) {
        err = domain.ErrClassroomStudentNotFound
        return
    }

    err = u.csRepo.Delete(ctx, id)
    if err != nil {
        return u.getError(err)
    }
    return
}

func (u *usecase) chunkSlice(slice []domain.ClassroomStudent, chunkSize int) [][]domain.ClassroomStudent {
	var chunks [][]domain.ClassroomStudent
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func (u *usecase) CopyClassroomStudent(c context.Context, classroomAcademicID int64, toClassroomAcademicID int64) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    ca, err := u.csaRepo.GetByID(ctx, classroomAcademicID)
    if err != nil {
        return u.getError(err)
    }

    if ca == (domain.ClassroomAcademic{}) {
        err = domain.ErrClassroomAcademicNotFound
        return
    }

    ca, err = u.csaRepo.GetByID(ctx, toClassroomAcademicID)
    if err != nil {
        return u.getError(err)
    }

    if ca == (domain.ClassroomAcademic{}) {
        err = domain.ErrClassroomAcademicNotFound
        return
    }

    res, err := u.GetByClassroomAcademic(ctx, classroomAcademicID)
    if err != nil {
        return 
    }

    var students []domain.ClassroomStudent
    for _, item := range res {
        student := domain.ClassroomStudent{
            ClassroomAcademic:  domain.ClassroomAcademic{
                ID: toClassroomAcademicID,
            },
            Student: item.Student,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }
        students = append(students, student)
    }

    chunk_students := u.chunkSlice(students, 100)
    for _, students := range chunk_students {
        err = u.csRepo.StoreMultiple(ctx, students)
        if err != nil {
            return u.getError(err)
        }
    }
    return
}