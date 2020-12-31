package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shellrean/extraordinary-raport/config"
	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/entities/DTO"
	"github.com/shellrean/extraordinary-raport/entities/helper"
	"github.com/shellrean/extraordinary-raport/interface/http/api"
	"github.com/shellrean/extraordinary-raport/interface/http/middleware"
)

type classHandler struct {
	classUsecase 		domain.ClassroomUsecase
	config				*config.Config
	mddl 				*middleware.GoMiddleware
}

func NewClassroomHandler(r *gin.Engine, m domain.ClassroomUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	handler := &classHandler {
		classUsecase:		m,
		config:				cfg,
		mddl:				mddl,
	}
	r.GET("/classrooms", handler.mddl.Auth(), handler.Fetch)
}

func (h *classHandler) Fetch(c *gin.Context) {
	res, err := h.classUsecase.Fetch(c)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}

	var data []dto.ClassroomResponse
	for _, item := range res {
		class := dto.ClassroomResponse {
			ID: 		item.ID,
			Name:		item.Name,
			MajorID: 	item.Major.ID,
			Grade:		item.Grade,
		}
		data = append(data, class)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
}