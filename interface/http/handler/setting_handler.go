package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/shellrean/extraordinary-raport/config"
    "github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/entities/DTO"
    "github.com/shellrean/extraordinary-raport/entities/helper"
    "github.com/shellrean/extraordinary-raport/interface/http/middleware"
    "github.com/shellrean/extraordinary-raport/interface/http/api"
)

type settingHanlder struct {
	settUsecase 		domain.SettingUsecase
	config 				*config.Config
	mddl 				*middleware.GoMiddleware
}

func NewSettingHandler(r *gin.Engine, m domain.SettingUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	handler := &settingHanlder {
		settUsecase:	m,
		config:			cfg,
		mddl:			mddl,	
	}

	setting := r.Group("/settings")
	setting.Use(handler.mddl.Auth())

	setting.GET("/", handler.Fetch)
}

func (h *settingHanlder) Fetch(c *gin.Context) {
	query, _ := c.GetQuery("q")
	if query == "" {
        err_code := helper.GetErrorCode(domain.ErrBadParamInput)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseError(domain.ErrBadParamInput.Error(), err_code),
        )
        return
	}
	
	nameV := strings.TrimRight(query, ",")
	nameV = strings.TrimLeft(nameV, ",") 
    names := strings.Split(nameV, ",")

    res, err := h.settUsecase.Fetch(c, names)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

	var data []dto.SettingResponse
	for _, item := range res {
		sett := dto.SettingResponse{
			Name:	item.Name,
			Value:	item.Value,
		}
		data = append(data, sett)
	}
    
    c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
}