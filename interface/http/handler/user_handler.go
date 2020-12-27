package handler

import (
    "net/http"
    "strings"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"

    "github.com/shellrean/extraordinary-raport/config"
    "github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/entities/helper"
    "github.com/shellrean/extraordinary-raport/interface/http/middleware"
    "github.com/shellrean/extraordinary-raport/interface/http/api"
)

type UserHandler struct {
    userUsecase     domain.UserUsecase
    config          *config.Config
    mddl            *middleware.GoMiddleware
}

func NewUserHandler(r *gin.Engine, m domain.UserUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
    handler := &UserHandler{
        userUsecase:    m,
        config:         cfg,
        mddl:           mddl,
    }
    r.GET("/users", handler.mddl.Auth(), handler.FetchUsers)
    r.GET("/users/:id", handler.mddl.Auth(), handler.Show)
    r.POST("/users", handler.mddl.Auth(), handler.Store)
    r.PUT("/users/:id", handler.mddl.Auth(), handler.Update)
    r.POST("/auth", handler.Autheticate)
    r.POST("/refresh-token", handler.RefreshToken)
}

func (h *UserHandler) FetchUsers(c *gin.Context) {
    limS , _ := c.GetQuery("limit")
    lim, _ := strconv.Atoi(limS)
    cursor, _ := c.GetQuery("cursor")

    res, nextCursor, err := h.userUsecase.Fetch(c, cursor, int64(lim))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    c.Header("X-Cursor", nextCursor)
    c.JSON(http.StatusOK, api.ResponseSuccess("success", res)) 
}

func (h *UserHandler) Autheticate(c *gin.Context) {
    var u domain.DTOUserLoginRequest
    if err := c.ShouldBindJSON(&u); err != nil {
        error_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(domain.ErrUnprocess.Error(), error_code),
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
        error_code := helper.GetErrorCode(domain.ErrValidation)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseErrorWithData(domain.ErrValidation.Error(), error_code, reserr),
        )
        return
    }

    res, err := h.userUsecase.Authentication(c, u)
    if err != nil {
        error_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), error_code),
        )
        return
    }
    
    c.JSON(http.StatusOK, api.ResponseSuccess("success",res))
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
    var u domain.DTOTokenResponse
    if err := c.ShouldBindJSON(&u); err != nil {
        error_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity, 
            api.ResponseError(domain.ErrUnprocess.Error(), error_code),
        )
        return
    }

    res, err := h.userUsecase.RefreshToken(c, u)
    if err != nil {
        error_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err), 
            api.ResponseError(err.Error(), error_code),
        )
        return
    }
    
    c.JSON(http.StatusOK, api.ResponseSuccess("success",res))
}

func (h *UserHandler) Show(c *gin.Context) {
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
    res := domain.DTOUserShow{}
    res, err = h.userUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("success", res))
}

func (h *UserHandler) Store(c *gin.Context) {
    var u domain.User
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

    res, err := h.userUsecase.Store(c, u)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("create user success", res))
}

func (h *UserHandler) Update(c *gin.Context) {
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
    res, err := h.userUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    if res == (domain.DTOUserShow{}) {
        err_code := helper.GetErrorCode(domain.ErrNotFound)
        c.JSON(
            api.GetHttpStatusCode(domain.ErrNotFound),
            api.ResponseError(domain.ErrNotFound.Error(), err_code),
        )
        return
    }

    var u domain.DTOUserUpdate
    err = c.ShouldBindJSON(&u)
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
    u.ID = int64(id)
    err = h.userUsecase.Update(c, &u)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("update user success", u))
}