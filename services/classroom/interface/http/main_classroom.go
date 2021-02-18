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
	"github.com/shellrean/extraordinary-raport/interface/http/api"
	"github.com/shellrean/extraordinary-raport/interface/http/middleware"
)

type handler struct {
	classUsecase 		domain.ClassroomUsecase
	config				*config.Config
	mddl 				*middleware.GoMiddleware
}

func New(r *gin.Engine, m domain.ClassroomUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	h := &handler {
		classUsecase:		m,
		config:				cfg,
		mddl:				mddl,
	}
	class := r.Group("/classrooms")
	class.Use(h.mddl.Auth())

	class.GET("/", h.Fetch)
	class.GET("/:id", h.Show)
	class.POST("/", h.Store)
	class.PUT("/:id", h.Update)
	class.DELETE("/:id", h.Delete)
}

func (h *handler) Fetch(c *gin.Context) {
	res, err := h.classUsecase.Fetch(c)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}

	var data []dto.ClassroomResponse
	for _, item := range res {
		class := dto.ClassroomResponse {
			ID: 		item.ID,
			Name:		item.Name,
			MajorID: 	item.Major.ID,
			Grade:		item.Grade,
		}
		data = append(data, class)
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
	res := domain.Classroom{}
	res, err = h.classUsecase.GetByID(c, int64(id))
	if err != nil {
		err_code := helper.GetErrorCode(err)
		c.JSON(
			http.StatusBadRequest,
			api.ResponseError(err.Error(), err_code),
		)
		return
	}
	
	data := dto.ClassroomResponse {
		ID:			res.ID,
		Name:		res.Name,
		Grade:		res.Grade,
		MajorID:	res.Major.ID,
	}
	c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}

func (h *handler) Store(c *gin.Context) {
	var cl dto.ClassroomRequest
	if err := c.ShouldBindJSON(&cl); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(domain.ErrUnprocess.Error(), err_code),
        )
        return
	}
	validate := validator.New()
    if err := validate.Struct(cl); err != nil {
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

	class := domain.Classroom {
		Name:		cl.Name,
		Grade:		cl.Grade,
	}
	major := domain.Major {
		ID:			cl.MajorID,
	}
	class.Major = major
	
	err := h.classUsecase.Store(c, &class)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseError(err.Error(), err_code),
        )
        return
	}

	data := dto.ClassroomResponse {
		ID:			class.ID,
		Name:		class.Name,
		Grade:		class.Grade,
		MajorID:	class.Major.ID,
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("create classroom success", data))
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
	res := domain.Classroom{}
	res, err = h.classUsecase.GetByID(c, int64(id))
	if err != nil {
		err_code := helper.GetErrorCode(err)
		c.JSON(
			http.StatusBadRequest,
			api.ResponseError(err.Error(), err_code),
		)
		return
	}

	if res == (domain.Classroom{}) {
		err_code := helper.GetErrorCode(domain.ErrNotFound)
		c.JSON(
			http.StatusBadRequest,
			api.ResponseError(domain.ErrNotFound.Error(), err_code),
		)
		return
	}

	var cl dto.ClassroomRequest
	if err := c.ShouldBindJSON(&cl); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(domain.ErrUnprocess.Error(), err_code),
        )
        return
	}
	validate := validator.New()
    if err := validate.Struct(cl); err != nil {
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

	class := domain.Classroom {
		ID:			int64(id),
		Name:		cl.Name,
		Grade:		cl.Grade,
	}
	major := domain.Major {
		ID:			cl.MajorID,
	}
	class.Major = major

	err = h.classUsecase.Update(c, &class)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseError(err.Error(), err_code),
        )
        return
	}

	data := dto.ClassroomResponse {
		ID:			class.ID,
		Name:		class.Name,
		Grade:		class.Grade,
		MajorID:	class.Major.ID,
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("update classroom success", data))
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
	res := domain.Classroom{}
	res, err = h.classUsecase.GetByID(c, int64(id))
	if err != nil {
		err_code := helper.GetErrorCode(err)
		c.JSON(
			http.StatusBadRequest,
			api.ResponseError(err.Error(), err_code),
		)
		return
	}

	if res == (domain.Classroom{}) {
		err_code := helper.GetErrorCode(domain.ErrNotFound)
		c.JSON(
			http.StatusBadRequest,
			api.ResponseError(domain.ErrNotFound.Error(), err_code),
		)
		return
	}

	err = h.classUsecase.Delete(c, int64(id))
	if err != nil {
		err_code := helper.GetErrorCode(err)
		c.JSON(
			http.StatusBadRequest,
			api.ResponseError(err.Error(), err_code),
		)
		return
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("delete classroom success", make([]string,0)))
}