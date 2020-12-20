package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/shellrean/extraordinary-raport/domain"
)

type UserHandler struct {
    userUsecase     domain.UserUsecase
}

func NewUserHandler(r *gin.Engine, m domain.UserUsecase) {
    handler := &UserHandler{
        userUsecase:    m,
    }
    r.GET("/users", handler.FetchUsers)
}

func (h *UserHandler) FetchUsers(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"user": "okeeee", "status": "no value"})
}
