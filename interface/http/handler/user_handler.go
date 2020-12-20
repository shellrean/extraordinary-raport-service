package handler

import (
    "net/http"
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"

    "github.com/shellrean/extraordinary-raport/config"
    "github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/entities/helper"
)

type UserHandler struct {
    userUsecase     domain.UserUsecase
    config          *config.Config
}

type ErrorResponse struct {
    Message     string  `json:"message"`
}

func NewUserHandler(r *gin.Engine, m domain.UserUsecase, cfg *config.Config) {
    handler := &UserHandler{
        userUsecase:    m,
        config:         cfg,
    }
    r.GET("/users", handler.FetchUsers)
    r.POST("/auth", handler.Autheticate)
}

func (h *UserHandler) FetchUsers(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"user": "okeeee", "status": "no value"}) 
}

func (h *UserHandler) Autheticate(c *gin.Context) {
    var u domain.DTOUserLoginRequest
    if err := c.ShouldBindJSON(&u); err != nil {
        c.JSON(http.StatusUnprocessableEntity, ErrorResponse{"Invalid json provided"})
        return
    }

    validate := validator.New()
    if err := validate.Struct(u); err != nil {
        var reserr []string

        errs := err.(validator.ValidationErrors)

		for _, e := range errs {
            res := fmt.Sprintf("%s error on tag %s",e.Field(), e.Tag())
            reserr = append(reserr, res)
        }
        c.JSON(http.StatusBadRequest, gin.H{"errors": reserr})
    }

    token, err := h.userUsecase.Authentication(c, u, h.config.JWTKey)
    if err != nil {
        c.JSON(helper.GetStatusCode(err), ErrorResponse{err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"token": token})
}