package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/shellrean/extraordinary-raport/config"
	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/entities/helper"
	"github.com/shellrean/extraordinary-raport/interface/http/api"
	"github.com/shellrean/extraordinary-raport/interface/http/middleware"
)

type studentHandler struct {
	studentUsecase		domain.StudentUsecase
	config				*config.Config
	mddl 				*middleware.GoMiddleware
}

func NewStudentHandler(r *gin.Engine, m domain.StudentUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	handler := &studentHandler{
		studentUsecase:	m,
		config: cfg,
		mddl:	mddl,
	}
    student := r.Group("/students")
    student.Use(handler.mddl.CORS())
    student.Use(handler.mddl.Auth())

	student.GET("/", handler.Index)
	student.POST("/", handler.Store)
	student.PUT("/:id", handler.Update)
}

func (h *studentHandler) Index(c *gin.Context) {
	limS, _ := c.GetQuery("limit")
	lim, _ := strconv.Atoi(limS)
	cursor, _ := c.GetQuery("cursor")

	res, nextCursor, err := h.studentUsecase.Fetch(c, cursor, int64(lim))
	if err != nil {
		error_code := helper.GetErrorCode(err)
		c.JSON(
			api.GetHttpStatusCode(err), 
			api.ResponseError(err.Error(), error_code),
		)
		return
	}

	c.Header("X-Cursor", nextCursor)
	c.JSON(http.StatusOK, api.ResponseSuccess("success",res))
}

func (h *studentHandler) Store(c *gin.Context) {
	var u domain.Student
    if err := c.ShouldBindJSON(&u); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(err.Error(), err_code),
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
	
	err := h.studentUsecase.Store(c, &u)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("create student success", u))
}

func (h *studentHandler) Update(c *gin.Context) {
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
    res, err := h.studentUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    if res == (domain.Student{}) {
        err_code := helper.GetErrorCode(domain.ErrNotFound)
        c.JSON(
            http.StatusNotFound,
            api.ResponseError(domain.ErrNotFound.Error(), err_code),
        )
        return
    }
	
	var u domain.Student
    if err := c.ShouldBindJSON(&u); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(err.Error(), err_code),
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
	u.ID = int64(id)
	err = h.studentUsecase.Update(c, &u)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("update student success", u))
}