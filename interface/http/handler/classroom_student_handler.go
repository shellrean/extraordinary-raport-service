package handler

import (
	"net/http"
	// "strings"
	"strconv"

	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"

	"github.com/shellrean/extraordinary-raport/config"
    "github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/entities/DTO"
    "github.com/shellrean/extraordinary-raport/entities/helper"
    "github.com/shellrean/extraordinary-raport/interface/http/middleware"
    "github.com/shellrean/extraordinary-raport/interface/http/api"
)

type csHandler struct {
	csUsecase 	domain.ClassroomStudentUsecase
	config 		*config.Config
	mddl		*middleware.GoMiddleware
}

func NewClassroomStudentHandler(
	r 	 *gin.Engine,
	m 	 domain.ClassroomStudentUsecase,
	cfg  *config.Config,
	mddl *middleware.GoMiddleware,
) {
	handler := &csHandler {
		csUsecase:		m,
		config:			cfg,
		mddl:			mddl,
	}
	cs := r.Group("/classroom-students")
	cs.Use(handler.mddl.Auth())

	cs.GET("/", handler.Fetch)
}

func (h *csHandler) Fetch(c *gin.Context) {
	limS , _ := c.GetQuery("limit")
    lim, _ := strconv.Atoi(limS)
	cursor, _ := c.GetQuery("cursor")
	
	res, nextCursor, err := h.csUsecase.Fetch(c, cursor, int64(lim))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	
	var data []dto.ClassroomStudentResponse
	for _, item := range res {
		ac := dto.ClassroomStudentResponse{
			ID:				item.ID,
			ClassroomAcademicID: item.ClassroomAcademic.ID,
			StudentID:		item.Student.ID,
		}
		data = append(data, ac)
	}

	c.Header("X-Cursor", nextCursor)
	c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
}
