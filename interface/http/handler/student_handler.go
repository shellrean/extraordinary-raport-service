package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"

	"github.com/shellrean/extraordinary-raport/config"
	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/entities/helper"
	"github.com/shellrean/extraordinary-raport/interface/http/api"
)

type studentHandler struct {
	studentUsecase		domain.StudentUsecase
	config				*config.Config
}

func NewStudentHandler(r *gin.Engine, m domain.StudentUsecase, cfg *config.Config) {
	handler := &studentHandler{
		studentUsecase:	m,
		config: cfg,
	}

	r.GET("/students", handler.Index)
}

func (h *studentHandler) Index(c *gin.Context) {
	res, err := h.studentUsecase.Fetch(c, int64(10))
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