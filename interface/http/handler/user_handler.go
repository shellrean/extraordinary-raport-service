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
)

type UserHandler struct {
    userUsecase     domain.UserUsecase
    config          *config.Config
    mddl            *middleware.GoMiddleware
}

type ErrorResponse struct {
    Message     string  `json:"message"`
}

type ErrorValidation struct {
    Field       string  `json:"field"`
    Message     string `json:"message"`
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
        c.JSON(http.StatusUnprocessableEntity, ErrorResponse{"invalid json provided"})
        return
    }

    validate := validator.New()
    if err := validate.Struct(u); err != nil {
        var reserr []ErrorValidation

        errs := err.(validator.ValidationErrors)

		for _, e := range errs {
            msg := helper.GetErrorMessage(e)
            
            res := ErrorValidation{
                Field:      strings.ToLower(e.Field()),
                Message:    msg,
            }

            reserr = append(reserr, res)
        }
        c.JSON(http.StatusBadRequest, gin.H{"message": "validation error","errors": reserr})
        return
    }

    res, err := h.userUsecase.Authentication(c, u)
    if err != nil {
        c.JSON(helper.GetStatusCode(err), ErrorResponse{err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, res)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
    var u domain.DTOTokenResponse
    if err := c.ShouldBindJSON(&u); err != nil {
        c.JSON(http.StatusUnprocessableEntity, ErrorResponse{"invalid json provided"})
        return
    }

    res, err := h.userUsecase.RefreshToken(c, u)
    if err != nil {
        c.JSON(helper.GetStatusCode(err), ErrorResponse{err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, res)
}