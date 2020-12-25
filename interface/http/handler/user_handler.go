package handler

import (
    "net/http"
    "strings"

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
    r.POST("/auth", handler.Autheticate)
    r.POST("/refresh-token", handler.RefreshToken)
}

func (h *UserHandler) FetchUsers(c *gin.Context) {
    data := c.GetStringMap("user-meta")
    c.JSON(http.StatusOK, gin.H{"data": data}) 
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