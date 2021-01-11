package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shellrean/extraordinary-raport/config"
	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/entities/DTO"
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
	academic := r.Group("/academics")
	academic.Use(handler.mddl.Auth())
	
	academic.GET("/", handler.Fetch)
	academic.POST("/", handler.Generate)
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
	var data []dto.AcademicResponse
	for _, item := range res {
		academic := dto.AcademicResponse {
			ID:		item.ID,
			Name: 	item.Name,
			Semester: item.Semester,
		}
		data = append(data, academic)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
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

	data := dto.AcademicResponse {
		ID: res.ID,
		Name: res.Name,
		Semester: res.Semester,
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success generate academic year",data)) 
}