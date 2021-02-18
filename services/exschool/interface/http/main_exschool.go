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

type handler struct {
	exschoolUsecase		domain.ExschoolUsecase
	config 				*config.Config
	mddl 				*middleware.GoMiddleware
}

func New(r *gin.Engine, m domain.ExschoolUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	h := &handler {
		exschoolUsecase: m,
		config: 		cfg,
		mddl: 			mddl,
	}

	ex := r.Group("/exschools")
	ex.Use(h.mddl.Auth())

	ex.GET("/", h.Fetch)
}

func (h *handler) Fetch(c *gin.Context) {
	res, err := h.exschoolUsecase.Fetch(c)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}


	var data []dto.ExschoolResponse
	for _, item := range res {
		exschool := dto.ExschoolResponse{
			ID: 	item.ID,
			Name: 	item.Name,
		}
		data = append(data, exschool)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}