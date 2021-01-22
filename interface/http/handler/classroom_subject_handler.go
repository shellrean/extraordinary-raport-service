package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/shellrean/extraordinary-raport/config"
    "github.com/shellrean/extraordinary-raport/domain"
    "github.com/shellrean/extraordinary-raport/entities/DTO"
    "github.com/shellrean/extraordinary-raport/entities/helper"
    "github.com/shellrean/extraordinary-raport/interface/http/middleware"
    "github.com/shellrean/extraordinary-raport/interface/http/api"
)

type csuHandler struct {
	csuUsecase 		domain.ClassroomSubjectUsecase
	config 			*config.Config
	mddl 			*middleware.GoMiddleware
}

func NewClassroomSubjectHandler(r *gin.Engine, m domain.ClassroomSubjectUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	handler := &csuHandler{
		csuUsecase:		m,
		config:			cfg,
		mddl:			mddl,
	}
	csu := r.Group("/classroom-subjects")
	csu.Use(handler.mddl.Auth())

	csu.GET("/classroom/:id", handler.FetchByClassroom)
}

func (h csuHandler) FetchByClassroom(c *gin.Context) {
	idS := c.Param("id")
    id, err := strconv.Atoi(idS)
    if err != nil {
        err_code := helper.GetErrorCode(domain.ErrBadParamInput)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseError(domain.ErrBadParamInput.Error(), err_code),
        )
        return
	}
	
	res, err := h.csuUsecase.FetchByClassroom(c, int64(id))
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}

	var data []dto.ClassroomSubjectResponse
	for _, item := range res {
		subject := dto.ClassroomSubjectResponse{
			ID:						item.ID,
			ClassroomAcademicID:	item.ClassroomAcademic.ID,
			SubjectID:				item.Subject.ID,
			SubjectName:			item.Subject.Name,
			TeacherID:				item.Teacher.ID,
			TeacherName:			item.Teacher.Name,
			MGN:					item.MGN,
		}

		data = append(data, subject)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}