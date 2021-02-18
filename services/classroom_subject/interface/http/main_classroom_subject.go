package handler

import (
	"net/http"
	"strings"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/shellrean/extraordinary-raport/config"
    "github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/entities/DTO"
    "github.com/shellrean/extraordinary-raport/entities/helper"
    "github.com/shellrean/extraordinary-raport/interface/http/middleware"
    "github.com/shellrean/extraordinary-raport/interface/http/api"
)

type handler struct {
	csuUsecase 		domain.ClassroomSubjectUsecase
	config 			*config.Config
	mddl 			*middleware.GoMiddleware
}

func New(r *gin.Engine, m domain.ClassroomSubjectUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	h := &handler{
		csuUsecase:		m,
		config:			cfg,
		mddl:			mddl,
	}
	csu := r.Group("/classroom-subjects")
	csu.Use(h.mddl.Auth())

    csu.GET("/", h.Fetch)
    csu.POST("/classroom-subject", h.Store)
    csu.GET("/classroom-subject/:id", h.Show)
    csu.PUT("/classroom-subject/:id", h.Update)
    csu.DELETE("/classroom-subject/:id", h.Delete)

    csu.GET("/classroom/:id", h.FetchByClassroom)
    csu.POST("/copy-subjects", h.CopyClassroomSubject)
}

func (h *handler) Fetch(c *gin.Context) {
    currentRoleUser := c.GetInt("role")
    currentUserID := c.GetInt64("user_id")
    user := domain.User{
        ID: currentUserID,
        Role: currentRoleUser,
    }
    
    res, err := h.csuUsecase.Fetch(c, user)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}

	var data []dto.ClassroomSubjectResponse
	for _, item := range res {
		subject := dto.ClassroomSubjectResponse{
			ID:						item.ID,
            ClassroomAcademicID:	item.ClassroomAcademic.ID,
            ClassroomName:          item.ClassroomAcademic.Classroom.Name,
			SubjectID:				item.Subject.ID,
			SubjectName:			item.Subject.Name,
			TeacherID:				item.Teacher.ID,
			TeacherName:			item.Teacher.Name,
			MGN:					item.MGN,
		}

		data = append(data, subject)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}

func (h *handler) Store(c *gin.Context) {
	var u dto.ClassroomSubjectRequest
    if err := c.ShouldBindJSON(&u); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(domain.ErrUnprocess.Error(), err_code),
        )
        return
    }
    validate := validator.New()
    if err := validate.Struct(u); err != nil {
        var reserr []api.ErrorValidation

        errs := err.(validator.ValidationErrors)
        for _, e := range errs {
            msg := helper.GetErrorMessage(e)
            res := api.ErrorValidation{
                Field:      strings.ToLower(e.Field()),
                Message:    msg,
            }
            reserr = append(reserr, res)
        }
        err_code := helper.GetErrorCode(domain.ErrValidation)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseErrorWithData(domain.ErrValidation.Error(), err_code, reserr),
        )
        return
	}
	
	classroomSubject := domain.ClassroomSubject{
		ClassroomAcademic:		domain.ClassroomAcademic{ID: u.ClassroomAcademicID},
		Subject:				domain.Subject{ID: u.SubjectID},
		Teacher:				domain.User{ID: u.TeacherID},
		MGN:					u.MGN,
	}

	err := h.csuUsecase.Store(c, &classroomSubject)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	
	u.ID = classroomSubject.ID
	c.JSON(http.StatusOK, api.ResponseSuccess("create classroom subject success", u))
}

func (h *handler) Update(c *gin.Context) {
    idS := c.Param("id")
    id, err := strconv.Atoi(idS)
    if err != nil {
        err_code := helper.GetErrorCode(domain.ErrBadParamInput)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseError(domain.ErrBadParamInput.Error(), err_code),
        )
        return
    }

	var u dto.ClassroomSubjectRequest
    if err := c.ShouldBindJSON(&u); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(domain.ErrUnprocess.Error(), err_code),
        )
        return
    }
    validate := validator.New()
    if err := validate.Struct(u); err != nil {
        var reserr []api.ErrorValidation

        errs := err.(validator.ValidationErrors)
        for _, e := range errs {
            msg := helper.GetErrorMessage(e)
            res := api.ErrorValidation{
                Field:      strings.ToLower(e.Field()),
                Message:    msg,
            }
            reserr = append(reserr, res)
        }
        err_code := helper.GetErrorCode(domain.ErrValidation)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseErrorWithData(domain.ErrValidation.Error(), err_code, reserr),
        )
        return
	}
	
	classroomSubject := domain.ClassroomSubject{
        ID:                     int64(id),
        ClassroomAcademic:		domain.ClassroomAcademic{ID: u.ClassroomAcademicID},
		Subject:				domain.Subject{ID: u.SubjectID},
		Teacher:				domain.User{ID: u.TeacherID},
		MGN:					u.MGN,
	}

	err = h.csuUsecase.Update(c, &classroomSubject)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	
	u.ID = classroomSubject.ID
	c.JSON(http.StatusOK, api.ResponseSuccess("update classroom subject success", u))
}

func (h *handler) Show(c *gin.Context) {
    idS := c.Param("id")
    id, err := strconv.Atoi(idS)
    if err != nil {
        err_code := helper.GetErrorCode(domain.ErrBadParamInput)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseError(domain.ErrBadParamInput.Error(), err_code),
        )
        return
    }
    
    res, err := h.csuUsecase.GetByID(c, int64(id))
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    
    data := dto.ClassroomSubjectResponse{
        ID:						res.ID,
        ClassroomAcademicID:	res.ClassroomAcademic.ID,
        ClassroomName:          res.ClassroomAcademic.Classroom.Name,
		SubjectID:				res.Subject.ID,
		SubjectName:			res.Subject.Name,
		TeacherID:				res.Teacher.ID,
		TeacherName:			res.Teacher.Name,
		MGN:					res.MGN,
    }

    c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}

func (h *handler) Delete(c *gin.Context) {
    idS := c.Param("id")
    id, err := strconv.Atoi(idS)
    if err != nil {
        err_code := helper.GetErrorCode(domain.ErrBadParamInput)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseError(domain.ErrBadParamInput.Error(), err_code),
        )
        return
    }

    err = h.csuUsecase.Delete(c, int64(id))
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    c.JSON(http.StatusOK, api.ResponseSuccess("success", make(map[string]string, 0)))
}

func (h *handler) FetchByClassroom(c *gin.Context) {
	idS := c.Param("id")
    id, err := strconv.Atoi(idS)
    if err != nil {
        err_code := helper.GetErrorCode(domain.ErrBadParamInput)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseError(domain.ErrBadParamInput.Error(), err_code),
        )
        return
	}
	
	res, err := h.csuUsecase.FetchByClassroom(c, int64(id))
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}

	var data []dto.ClassroomSubjectResponse
	for _, item := range res {
		subject := dto.ClassroomSubjectResponse{
			ID:						item.ID,
            ClassroomAcademicID:	item.ClassroomAcademic.ID,
            ClassroomName:          item.ClassroomAcademic.Classroom.Name,
			SubjectID:				item.Subject.ID,
			SubjectName:			item.Subject.Name,
			TeacherID:				item.Teacher.ID,
			TeacherName:			item.Teacher.Name,
			MGN:					item.MGN,
		}

		data = append(data, subject)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}

func (h *handler) CopyClassroomSubject(c *gin.Context) {
    var u dto.ClassroomSubjectCopyRequest
    if err := c.ShouldBindJSON(&u); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(domain.ErrUnprocess.Error(), err_code),
        )
        return
    }

    err := h.csuUsecase.CopyClassroomSubject(c, u.ClassroomAcademicID, u.ToClassroomAcademicID)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	
	c.JSON(http.StatusOK, api.ResponseSuccess("copy classroom subject success", make(map[string]string, 0)))
}