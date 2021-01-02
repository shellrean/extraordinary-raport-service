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

type classAcademicHandler struct {
	classAcademicUsecase domain.ClassroomAcademicUsecase
	config 				 *config.Config
	mddl 				 *middleware.GoMiddleware
}

func NewClassAcademicHandler(
	r 		*gin.Engine,
	m 		domain.ClassroomAcademicUsecase,
	cfg 	*config.Config,
	mddl 	*middleware.GoMiddleware,
) {
	handler := &classAcademicHandler {
		classAcademicUsecase:		m,
		config:						cfg,
		mddl:						mddl,
	}
	ca := r.Group("/classroom-academics")
	ca.Use(handler.mddl.CORS())
	ca.Use(handler.mddl.Auth())

	ca.GET("/", handler.Fetch)
}	

func (h *classAcademicHandler) Fetch(c *gin.Context) {
	res, err := h.classAcademicUsecase.Fetch(c)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	
	var data []dto.ClassroomAcademicResponse
	for _, item := range res {
		ac := dto.ClassroomAcademicResponse{
			ID:				item.ID,
			AcademicID:		item.Academic.ID,
			TeacherID:		item.Teacher.ID,
		}
		data = append(data, ac)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
}