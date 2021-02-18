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
	csUsecase 	domain.ClassroomStudentUsecase
	config 		*config.Config
	mddl		*middleware.GoMiddleware
}

func New(
	r 	 *gin.Engine,
	m 	 domain.ClassroomStudentUsecase,
	cfg  *config.Config,
	mddl *middleware.GoMiddleware,
) {
	h := &handler {
		csUsecase:		m,
		config:			cfg,
		mddl:			mddl,
	}
	cs := r.Group("/classroom-students")
	cs.Use(h.mddl.Auth())

	cs.GET("/", h.Fetch)
    cs.GET("/classroom-student/:id", h.Show)
    cs.POST("/classroom-student", h.Store)
    cs.PUT("/classroom-student/:id", h.Update)
    cs.DELETE("/classroom-student/:id", h.Delete)
    
    cs.GET("/classroom/:id", h.FetchByClassroom)
    cs.POST("/copy-students", h.CopyClassroomStudent)
}

func (h *handler) Fetch(c *gin.Context) {
	limS , _ := c.GetQuery("limit")
    lim, _ := strconv.Atoi(limS)
	cursor, _ := c.GetQuery("cursor")
	
	res, nextCursor, err := h.csUsecase.Fetch(c, cursor, int64(lim))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	
	var data []dto.ClassroomStudentResponse
	for _, item := range res {
		ac := dto.ClassroomStudentResponse{
			ID:				item.ID,
            ClassroomAcademicID: item.ClassroomAcademic.ID,
            StudentID:		item.Student.ID,
            StudentSRN:     item.Student.SRN,
            StudentNSRN:    item.Student.NSRN,
            StudentName:    item.Student.Name,
		}
		data = append(data, ac)
	}

	c.Header("X-Cursor", nextCursor)
	c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
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
	
	res, err := h.csUsecase.GetByClassroomAcademic(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	
	var data []dto.ClassroomStudentResponse
	for _, item := range res {
		ac := dto.ClassroomStudentResponse{
			ID:				item.ID,
            ClassroomAcademicID: item.ClassroomAcademic.ID,
            StudentID:		item.Student.ID,
            StudentSRN:     item.Student.SRN,
            StudentNSRN:    item.Student.NSRN,
            StudentName:    item.Student.Name,
		}
		data = append(data, ac)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
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
    res := domain.ClassroomStudent{}
    res, err = h.csUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    data := dto.ClassroomStudentResponse {
        ID:				res.ID,
		ClassroomAcademicID: res.ClassroomAcademic.ID,
        StudentID:		res.Student.ID,
        StudentSRN:     res.Student.SRN,
        StudentNSRN:    res.Student.NSRN,
        StudentName:    res.Student.Name,
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}

func (h *handler) Store(c *gin.Context) {
    var u dto.ClassroomStudentRequest
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

    cs := domain.ClassroomStudent{
        ClassroomAcademic:   domain.ClassroomAcademic{ID: u.ClassroomAcademicID},
        Student:   			 domain.Student{ID: u.StudentID},
    }

    err := h.csUsecase.Store(c, &cs)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    u.ID = cs.ID
    c.JSON(http.StatusOK, api.ResponseSuccess("create classroom student success", u))
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
    res, err := h.csUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    if res == (domain.ClassroomStudent{}) {
        err_code := helper.GetErrorCode(domain.ErrNotFound)
        c.JSON(
            api.GetHttpStatusCode(domain.ErrNotFound),
            api.ResponseError(domain.ErrNotFound.Error(), err_code),
        )
        return
    }

    var u dto.ClassroomStudentRequest
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

    cs := domain.ClassroomStudent{
        ID:                  int64(id),
        ClassroomAcademic:   domain.ClassroomAcademic{ID: u.ClassroomAcademicID},
        Student:   			 domain.Student{ID: u.StudentID},
    }

    err = h.csUsecase.Update(c, &cs)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    u.ID = cs.ID
    c.JSON(http.StatusOK, api.ResponseSuccess("update classroom's student success", u))
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
    res, err := h.csUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    if res == (domain.ClassroomStudent{}) {
        err_code := helper.GetErrorCode(domain.ErrNotFound)
        c.JSON(
            api.GetHttpStatusCode(domain.ErrNotFound),
            api.ResponseError(domain.ErrNotFound.Error(), err_code),
        )
        return
    }

    err = h.csUsecase.Delete(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("delete classroom's student success", make([]string,0)))
}

func (h *handler) CopyClassroomStudent(c *gin.Context) {
    var u dto.ClassroomStudentCopyRequest
    if err := c.ShouldBindJSON(&u); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(domain.ErrUnprocess.Error(), err_code),
        )
        return
    }

    err := h.csUsecase.CopyClassroomStudent(c, u.ClassroomAcademicID, u.ToClassroomAcademicID)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("copy classroom's student success", make([]string,0)))
}