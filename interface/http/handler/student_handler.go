package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"

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

	r.GET("/students", handler.mddl.Auth(), handler.Index)
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