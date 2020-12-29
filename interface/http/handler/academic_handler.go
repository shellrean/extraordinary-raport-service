package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"

	"github.com/shellrean/extraordinary-raport/config"
    "github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/entities/helper"
    "github.com/shellrean/extraordinary-raport/interface/http/middleware"
    "github.com/shellrean/extraordinary-raport/interface/http/api"
)

type academicHandler struct {
	academicUsecase 		domain.AcademicUsecase
	config 					*config.Config
	mddl 					*middleware.GoMiddleware
}

func NewAcademicHandler(r *gin.Engine, m domain.AcademicUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	handler := &academicHandler {
		academicUsecase:		m,
		config:					cfg,
		mddl:					mddl,
	}
	r.GET("/academics", handler.mddl.Auth(), handler.Fetch)
	r.GET("/academics-generate", handler.mddl.Auth(), handler.Generate)
}

func (h *academicHandler) Fetch(c *gin.Context) {
	res, err := h.academicUsecase.Fetch(c)
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

func (h *academicHandler) Generate(c *gin.Context) {
	res, err := h.academicUsecase.Generate(c)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	c.JSON(http.StatusOK, api.ResponseSuccess("success generate academic year", res)) 
}