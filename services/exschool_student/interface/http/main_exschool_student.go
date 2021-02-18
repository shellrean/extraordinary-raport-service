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
	exsUsecase 		domain.ExschoolStudentUsecase
	config 			*config.Config
	mddl 			*middleware.GoMiddleware
}

func New(r *gin.Engine, m domain.ExschoolStudentUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	h := &handler{
		exsUsecase:		m,
		config:			cfg,
		mddl:			mddl,
	}

	exs := r.Group("/exschool-students")
	exs.Use(h.mddl.Auth())

	exs.GET("/", h.Fetch)
	exs.POST("exschool-student", h.Store)
	exs.DELETE("exschool-student/:id", h.Delete)

	exs.GET("/classroom/:id", h.FetchByClassroom)
}

func (h *handler) Fetch(c *gin.Context) {
	return
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
	
	res, err := h.exsUsecase.FetchByClassroom(c, int64(id))
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}

	var data []dto.ExschoolStudentResponse
	for _, item := range res {
		exs := dto.ExschoolStudentResponse{
			ID:			item.ID,
			ExschoolID: item.Exschool.ID,
			StudentID:	item.Student.ID,
		}
		data = append(data, exs)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}

func (h *handler) Store(c *gin.Context) {
	var u dto.ExschoolStudentRequest
    err := c.ShouldBindJSON(&u)
    if err != nil {
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
	
	exs := domain.ExschoolStudent {
		Exschool:	domain.Exschool{
			ID:	 	u.ExschoolID,
		},
		Student:	domain.ClassroomStudent{
			ID: 	u.StudentID,
		},
	}

	err = h.exsUsecase.Store(c, &exs)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	
	data := u
	data.ID = exs.ID
	c.JSON(http.StatusOK, api.ResponseSuccess("create exschool student success", data))
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
	
	err = h.exsUsecase.Delete(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("delete exschool student success", make([]string,0)))
}