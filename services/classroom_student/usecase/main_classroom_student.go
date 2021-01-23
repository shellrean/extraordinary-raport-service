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

type csUsecase struct {
    csRepo			domain.ClassroomStudentRepository
    csaRepo         domain.ClassroomAcademicRepository
    settingRepo     domain.SettingRepository
	contextTimeout	time.Duration
	cfg 			*config.Config
}

func NewClassroomStudentUsecase(
    d domain.ClassroomStudentRepository,
    m domain.ClassroomAcademicRepository,
    sr domain.SettingRepository,
	timeout time.Duration,
	cfg *config.Config,
) domain.ClassroomStudentUsecase {
	return &csUsecase {
        csRepo:			d,
        csaRepo:        m,
        settingRepo:    sr,
		contextTimeout: timeout,
		cfg:			cfg,
	}
}

func (u *csUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.ClassroomStudent, nextCursor string, err error) {
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
        if u.cfg.Release {
            log.Println(err.Error())
            err = domain.ErrServerError
            return
        }
        return
    }

    if len(res) == int(num) {
        nextCursor = helper.EncodeCursor(res[len(res)-1].ID)
    }

    return
}

func (u *csUsecase) GetByClassroomAcademic(c context.Context, classroomAcademicID int64) (res []domain.ClassroomStudent, err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    csa, err := u.csaRepo.GetByID(ctx, classroomAcademicID)
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

    res, err = u.csRepo.GetByClassroomAcademic(ctx, classroomAcademicID)
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

func (u *csUsecase) GetByID(c context.Context, id int64) (res domain.ClassroomStudent, err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    res, err = u.csRepo.GetByID(ctx, id)
    if err != nil {
        if u.cfg.Release {
            log.Println(err.Error())
            err = domain.ErrServerError
            return
        }
        return 
    }
    if res == (domain.ClassroomStudent{}) {
        err = domain.ErrClassroomStudentNotFound
        return
    }
    return
}

func (u *csUsecase) Store(c context.Context, cs *domain.ClassroomStudent) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    res, err := u.settingRepo.GetByName(ctx, domain.SettingAcademicActive)
	if err != nil {
		if u.cfg.Release {
            log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
    }
    
    if res == (domain.Setting{}) {
		err = domain.ErrSettingNotFound
		return
    }
    
    id, err := strconv.Atoi(res.Value)
	if err != nil {
		if u.cfg.Release {
            log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
    }
    
    exist, err := u.csRepo.GetByAcademicAndStudent(ctx, int64(id), cs.Student.ID)
	if err != nil {
		if u.cfg.Release {
            log.Println(err.Error())
			err = domain.ErrServerError
			return
		}
		return
	}

	if exist != (domain.ClassroomStudent{}) {
        err = domain.ErrClassroomStudentExist
        return
	}

    cs.CreatedAt = time.Now()
    cs.UpdatedAt = time.Now()

    if err = u.csRepo.Store(ctx, cs); err != nil {
        if u.cfg.Release {
            log.Println(err.Error())
            err = domain.ErrServerError
            return
        }
        return
    }

    return
}

func (u *csUsecase) Update(c context.Context, cs *domain.ClassroomStudent) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    exist, err := u.csRepo.GetByID(ctx, cs.ID)
    if err != nil {
        if u.cfg.Release {
            log.Println(err.Error())
            return domain.ErrServerError
        }
        return
    }

    if exist == (domain.ClassroomStudent{}) {
        err = domain.ErrClassroomStudentNotFound
        return
    }

    cs.UpdatedAt = time.Now()

    if err = u.csRepo.Update(ctx, cs); err != nil {
        if u.cfg.Release {
            log.Println(err.Error())
            return domain.ErrServerError
        }
        return
    }

    return
}

func (u *csUsecase) Delete(c context.Context, id int64) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    exist, err := u.csRepo.GetByID(ctx, id)
    if err != nil {
        if u.cfg.Release {
            log.Println(err.Error())
            return domain.ErrServerError
        }
        return
    }

    if exist == (domain.ClassroomStudent{}) {
        err = domain.ErrClassroomStudentNotFound
        return
    }

    err = u.csRepo.Delete(ctx, id)
    if err != nil {
        if u.cfg.Release {
            log.Println(err.Error())
            return domain.ErrServerError
        }
        return
    }
    return
}

func (u *csUsecase) chunkSlice(slice []domain.ClassroomStudent, chunkSize int) [][]domain.ClassroomStudent {
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

func (u *csUsecase) CopyClassroomStudent(c context.Context, classroomAcademicID int64, toClassroomAcademicID int64) (err error) {
    ctx, cancel := context.WithTimeout(c, u.contextTimeout)
    defer cancel()

    ca, err := u.csaRepo.GetByID(ctx, classroomAcademicID)
    if err != nil {
        if u.cfg.Release {
            log.Println(err.Error())
            err = domain.ErrServerError
            return
        }
        return
    }

    if ca == (domain.ClassroomAcademic{}) {
        err = domain.ErrClassroomAcademicNotFound
        return
    }

    ca, err = u.csaRepo.GetByID(ctx, toClassroomAcademicID)
    if err != nil {
        if u.cfg.Release {
            log.Println(err.Error())
            err = domain.ErrServerError
            return
        }
        return
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
            if u.cfg.Release {
                log.Println(err.Error())
                err = domain.ErrServerError
                return
            }
            return
        }
    }
    return
}