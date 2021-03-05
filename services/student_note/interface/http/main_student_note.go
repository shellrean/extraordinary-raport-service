package handler

import (
	"net/http"
	"strings"
	"strconv"
	
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
	snUsecase 		domain.StudentNoteUsecase
	config 			*config.Config
	mddl 			*middleware.GoMiddleware
}

func New(r *gin.Engine, m domain.StudentNoteUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	h := &handler {
		snUsecase: 	m,
		config:		cfg,
		mddl:		mddl,
	}

	sn := r.Group("/student-notes")
	sn.Use(h.mddl.Auth())

	sn.POST("student-note",h.mddl.MustRole([]int{domain.RoleTeacher}), h.Store)
	sn.GET("classroom/:id", h.mddl.MustRole([]int{domain.RoleTeacher}), h.FetchByClassroom)
}

func (h *handler) FetchByClassroom(c *gin.Context) {
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

	typ := c.Query("type")
	var typID int64
	if typ != "" {
		typID, err = strconv.ParseInt(typ, 10, 64)
		if err != nil {
			err_code := helper.GetErrorCode(domain.ErrBadParamInput)
			c.JSON(
				http.StatusBadRequest,
				api.ResponseError(domain.ErrBadParamInput.Error(), err_code),
			)
			return
		}
	}

	res, err := h.snUsecase.FetchByClassroom(c, int64(id), typID)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}

	var data []dto.StudentNoteResponse
	for _, item := range res {
		snr := dto.StudentNoteResponse{
			ID: 		item.ID,
			Type:		item.Type,
			StudentID:	item.Student.ID,
			TeacherID:	item.Teacher.ID,
			Note:		item.Note,
		}
		data = append(data, snr)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}

func (h *handler) Store(c *gin.Context) {
	var u dto.StudentNoteRequest
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
    currentUserID := c.GetInt64("user_id")
    u.TeacherID = currentUserID
	
	sn := domain.StudentNote {
		Type:		u.Type,
		Teacher: 	domain.User{
			ID: 	u.TeacherID,
		},
		Student: 	domain.ClassroomStudent{
			ID:		u.StudentID,
		},
		Note:		u.Note,
	}
	err := h.snUsecase.Store(c, &sn)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}

	u.ID = sn.ID
	c.JSON(http.StatusOK, api.ResponseSuccess("create student note success", u))
}