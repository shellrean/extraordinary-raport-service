package handler

import (
	"net/http"
    "strconv"
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

type handler struct {
    attendUsecase   domain.AttendanceUsecase
    config          *config.Config
    mddl            *middleware.GoMiddleware
}

func New(r *gin.Engine, m domain.AttendanceUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	h := &handler{
        attendUsecase:  m,
        config:         cfg,
        mddl:           mddl,
	}
	attend := r.Group("attendances")
	attend.Use(h.mddl.Auth())

    attend.GET("/:id", h.Fetch)
    attend.POST("/", h.mddl.MustRole([]int{domain.RoleTeacher}), h.Store)
}

func (h *handler) Fetch(c *gin.Context) {
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
	
	res, err := h.attendUsecase.Fetch(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	
	var data []dto.AttendanceResponse
	for _, item := range res {
		attendance := dto.AttendanceResponse{
			ID:			item.ID,
            StudentID:	item.Student.ID,
            Type:       item.Type,
			Total:		item.Total,
		}
		data = append(data, attendance)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
}

func (h *handler) Store(c *gin.Context) {
    var u dto.AttendanceRequest
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

    attener := domain.Attendance{
        Student:        domain.ClassroomStudent{
            ID:     u.StudentID,
        },
        Total:          u.Total,
        Type:           u.Type,
    }
    
    err := h.attendUsecase.Store(c, &attener)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    data := dto.AttendanceResponse{
        ID:         attener.ID,
        StudentID:  attener.Student.ID,
        Total:      attener.Total,
        Type:       attener.Type,
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("create attendance success", data))
}