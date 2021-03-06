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
	classAcademicUsecase domain.ClassroomAcademicUsecase
	config 				 *config.Config
	mddl 				 *middleware.GoMiddleware
}

func New(
	r 		*gin.Engine,
	m 		domain.ClassroomAcademicUsecase,
	cfg 	*config.Config,
	mddl 	*middleware.GoMiddleware,
) {
	h := &handler {
		classAcademicUsecase:		m,
		config:						cfg,
		mddl:						mddl,
	}
	ca := r.Group("/classroom-academics")
	ca.Use(h.mddl.Auth())

    ca.GET("/", h.Fetch)
    ca.GET("/classroom-academic/:id", h.Show)
	ca.POST("/classroom-academic", h.Store)
	ca.PUT("/classroom-academic/:id", h.Update)
    ca.DELETE("/classroom-academic/:id", h.Delete)
    
    ca.GET("/academic/:id", h.FetchByAcademic)
}	

func (h *handler) Fetch(c *gin.Context) {
    currentRoleUser := c.GetInt("role")
    currentUserID := c.GetInt64("user_id")
    user := domain.User{
        ID: currentUserID,
        Role: currentRoleUser,
    }
    
    res, err := h.classAcademicUsecase.Fetch(c, user)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
	
	var data []dto.ClassroomAcademicResponse
	for _, item := range res {
		ac := dto.ClassroomAcademicResponse{
			ID:				item.ID,
			AcademicID:		item.Academic.ID,
			TeacherID:		item.Teacher.ID,
            ClassroomID: 	item.Classroom.ID,
            TeacherName:    item.Teacher.Name,
            ClassroomName:  item.Classroom.Name,
            ClassroomMajor: item.Classroom.Major.Name,
		}
		data = append(data, ac)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
}

func (h *handler) FetchByAcademic(c *gin.Context) {
    idS := c.Param("id")
	AcademicID, err := strconv.Atoi(idS)
	if err != nil {
		err_code := helper.GetErrorCode(domain.ErrBadParamInput)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseError(domain.ErrBadParamInput.Error(), err_code),
        )
        return
    }

    res, err := h.classAcademicUsecase.FetchByAcademic(c, int64(AcademicID))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	
	var data []dto.ClassroomAcademicResponse
	for _, item := range res {
		ac := dto.ClassroomAcademicResponse{
			ID:				item.ID,
			AcademicID:		item.Academic.ID,
			TeacherID:		item.Teacher.ID,
            ClassroomID: 	item.Classroom.ID,
            TeacherName:    item.Teacher.Name,
            ClassroomName:  item.Classroom.Name,
            ClassroomMajor: item.Classroom.Major.Name,
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
    
	res, err := h.classAcademicUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	
    data := dto.ClassroomAcademicResponse{
		ID:				res.ID,
		AcademicID:		res.Academic.ID,
		TeacherID:		res.Teacher.ID,
        ClassroomID: 	res.Classroom.ID,
        TeacherName:    res.Teacher.Name,
        ClassroomName:  res.Classroom.Name,
        ClassroomMajor: res.Classroom.Major.Name,
    }
    
	c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
}

func (h *handler) Store(c *gin.Context) {
	var u dto.ClassroomAcademicRequest
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

	ca := domain.ClassroomAcademic {
		Teacher:	domain.User {
			ID:		u.TeacherID,
		},
		Classroom: 	domain.Classroom {
			ID: 	u.ClassroomID,
		},
	}

	err := h.classAcademicUsecase.Store(c, &ca)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}


	u.ID = ca.ID
	c.JSON(http.StatusOK, api.ResponseSuccess("success create academic classroom", u))
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
	res, err := h.classAcademicUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    if res == (domain.ClassroomAcademic{}) {
        err_code := helper.GetErrorCode(domain.ErrNotFound)
        c.JSON(
            api.GetHttpStatusCode(domain.ErrNotFound),
            api.ResponseError(domain.ErrNotFound.Error(), err_code),
        )
        return
	}
	
	var u dto.ClassroomAcademicRequest
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

	ca := domain.ClassroomAcademic {
		ID:	int64(id),
		Teacher:	domain.User {
			ID:		u.TeacherID,
		},
		Classroom: 	domain.Classroom {
			ID: 	u.ClassroomID,
		},
	}
	err = h.classAcademicUsecase.Update(c, &ca)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	u.ID = int64(id)
	c.JSON(http.StatusOK, api.ResponseSuccess("success update academic classroom", u))
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

	res, err := h.classAcademicUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    if res == (domain.ClassroomAcademic{}) {
        err_code := helper.GetErrorCode(domain.ErrNotFound)
        c.JSON(
            api.GetHttpStatusCode(domain.ErrNotFound),
            api.ResponseError(domain.ErrNotFound.Error(), err_code),
        )
        return
	}

	err = h.classAcademicUsecase.Delete(c, int64(id))
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("delete classroom academic success", make([]string,0)))
}