package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

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
	ca.POST("/", handler.Store)
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
			ClassroomID: 	item.Classroom.ID,
		}
		data = append(data, ac)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
}

func (h *classAcademicHandler) Store(c *gin.Context) {
	var u dto.ClassroomAcademicResponse
    if err := c.ShouldBindJSON(&u); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(domain.ErrUnprocess.Error(), err_code),
        )
        return
	}
	validate := validator.New()
    if err := validate.Struct(u); err != nil {
        var reserr []api.ErrorValidation

        errs := err.(validator.ValidationErrors)
        for _, e := range errs {
            msg := helper.GetErrorMessage(e)
            res := api.ErrorValidation{
                Field:      strings.ToLower(e.Field()),
                Message:    msg,
            }
            reserr = append(reserr, res)
        }
        err_code := helper.GetErrorCode(domain.ErrValidation)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseErrorWithData(domain.ErrValidation.Error(), err_code, reserr),
        )
        return
	}

	ca := domain.ClassroomAcademic {
		Academic: 	domain.Academic {
			ID: 	u.AcademicID,
		},
		Teacher:	domain.User {
			ID:		u.TeacherID,
		},
		Classroom: 	domain.Classroom {
			ID: 	u.ClassroomID,
		},
	}

	err := h.classAcademicUsecase.Store(c, &ca)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseError(err.Error(), err_code),
        )
        return
	}


	u.ID = ca.ID
	c.JSON(http.StatusOK, api.ResponseSuccess("success create academic classroom", u))
}