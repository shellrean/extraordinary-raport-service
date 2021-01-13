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

type subjectHandler struct {
    subjectUsecase  domain.SubjectUsecase
    config          *config.Config
    mddl            *middleware.GoMiddleware
}

func NewSubjectHandler(r *gin.Engine, m domain.SubjectUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	handler := &subjectHandler{
		subjectUsecase: 	m,
		config:				cfg,
		mddl:				mddl,
	}

	subject := r.Group("/subjects")
	subject.Use(handler.mddl.Auth())

    subject.GET("/", handler.Fetch)
    subject.GET("/:id", handler.Show)
    subject.POST("/", handler.Store)
    subject.PUT("/:id", handler.Update)
    subject.DELETE("/:id", handler.Delete)
}

func (h *subjectHandler) Fetch(c *gin.Context) {
	limS , _ := c.GetQuery("limit")
    lim, _ := strconv.Atoi(limS)
    cursor, _ := c.GetQuery("cursor")

    res, nextCursor, err := h.subjectUsecase.Fetch(c, cursor, int64(lim))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    var data []dto.SubjectResponse
    for _, item := range res {
        user := dto.SubjectResponse{
            ID:     item.ID,
			Name:   item.Name,
			Type:	item.Type,
        }
        data = append(data, user)
    }

    c.Header("X-Cursor", nextCursor)
    c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
}

func (h *subjectHandler) Show(c *gin.Context) {
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
    res := domain.Subject{}
    res, err = h.subjectUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    data := dto.SubjectResponse {
        ID: res.ID,
        Name: res.Name,
        Type: res.Type,
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}


func (h *subjectHandler) Store(c *gin.Context) {
    var u dto.SubjectResponse
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

    subject := domain.Subject{
        Name:   u.Name,
        Type:   u.Type,
    }

    err := h.subjectUsecase.Store(c, &subject)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    u.ID = subject.ID
    c.JSON(http.StatusOK, api.ResponseSuccess("create subject success", u))
}

func (h *subjectHandler) Update(c *gin.Context) {
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
    res, err := h.subjectUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    if res == (domain.Subject{}) {
        err_code := helper.GetErrorCode(domain.ErrNotFound)
        c.JSON(
            api.GetHttpStatusCode(domain.ErrNotFound),
            api.ResponseError(domain.ErrNotFound.Error(), err_code),
        )
        return
    }

    var u dto.SubjectResponse
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
    u.ID = int64(id)

    subject := domain.Subject{
        ID:     u.ID,
        Name:   u.Name,
        Type:   u.Type,
    }

    err = h.subjectUsecase.Update(c, &subject)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("update subject success", u))
}

func (h *subjectHandler) Delete(c *gin.Context) {
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
    res, err := h.subjectUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    if res == (domain.Subject{}) {
        err_code := helper.GetErrorCode(domain.ErrNotFound)
        c.JSON(
            api.GetHttpStatusCode(domain.ErrNotFound),
            api.ResponseError(domain.ErrNotFound.Error(), err_code),
        )
        return
    }

    err = h.subjectUsecase.Delete(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("delete subject success", make([]string,0)))
}