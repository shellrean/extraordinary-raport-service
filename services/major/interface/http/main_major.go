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
    majorUsecase     domain.MajorUsecase
    config          *config.Config
    mddl            *middleware.GoMiddleware
}

func New(r *gin.Engine, m domain.MajorUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	h := &handler{
        majorUsecase:    m,
        config:         cfg,
        mddl:           mddl,
	}
	major := r.Group("/majors")
	major.Use(h.mddl.Auth())

	major.GET("/", h.Fetch)
}

func (h *handler) Fetch(c *gin.Context) {
	res, err := h.majorUsecase.Fetch(c)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	
	var data []dto.MajorResponse
	for _, item := range res {
		major := dto.MajorResponse{
			ID:			item.ID,
			Name:		item.Name,
		}
		data = append(data, major)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
}