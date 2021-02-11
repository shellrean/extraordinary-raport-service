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

type cspHandler struct {
	cspUsecase	domain.ClassroomSubjectPlanUsecase
	config 		*config.Config
	mddl 		*middleware.GoMiddleware
}

func NewClassroomSubjectPlanHandler(r *gin.Engine, m domain.ClassroomSubjectPlanUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	handler := &cspHandler {
		cspUsecase:	m,
		config:		cfg,
		mddl:		mddl,
	}
	csp := r.Group("/classroom-subject-plans")
	csp.Use(handler.mddl.Auth())

    csp.POST("/", handler.mddl.MustRole([]int{domain.RoleTeacher, domain.RoleAdmin}), handler.Fetch)
	csp.POST("/classroom-subject-plan", handler.mddl.MustRole([]int{domain.RoleTeacher}), handler.Store)
	csp.PUT("/classroom-subject-plan", handler.mddl.MustRole([]int{domain.RoleTeacher}), handler.Update)
	csp.DELETE("/classroom-subject-plan/:id", handler.mddl.MustRole([]int{domain.RoleTeacher}), handler.Delete)
}

func (h cspHandler) Fetch(c *gin.Context) {
    var u dto.ClassroomSubjectPlanFetchRequest
    if err := c.ShouldBindJSON(&u); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(domain.ErrUnprocess.Error(), err_code),
        )
        return
    }

    currentRoleUser := c.GetInt("role")
    currentUserID := c.GetInt64("user_id")

    if currentRoleUser == domain.RoleTeacher {
        u.TeacherID = currentUserID
    }

    res, err := h.cspUsecase.Fetch(c, u.Query, u.TeacherID, u.ClassroomID)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    var data []dto.ClassroomSubjectPlanResponse
    for _, item := range res {
        csp := dto.ClassroomSubjectPlanResponse {
            ID: item.ID,
            Type: item.Type,
            Name: item.Name,
            Desc: item.Desc,
            TeacherID: item.Teacher.ID,
            SubjectID: item.Subject.ID,
            SubjectName: item.Subject.Subject.Name,
            ClassroomID: item.Classroom.ID,
            CountPlan: item.CountPlan,
            MaxPoint: item.MaxPoint,
        }

        data = append(data, csp)
    }

    c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}

func (h cspHandler) Store(c *gin.Context) {
	var u dto.ClassroomSubjectPlanRequest
	if err := c.ShouldBindJSON(&u); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(domain.ErrUnprocess.Error(), err_code),
        )
        return
    }
    
    currentRoleUser := c.GetInt("role")
    currentUserID := c.GetInt64("user_id")

    if currentRoleUser == domain.RoleTeacher {
        u.TeacherID = currentUserID
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

	csp := domain.ClassroomSubjectPlan {
		Type: 			u.Type,
		Name:			u.Name,
		Desc:			u.Desc,
		Teacher:		domain.User{
			ID:	u.TeacherID,
		},
		Subject:		domain.ClassroomSubject{
			ID: u.SubjectID,
		},
		Classroom: 		domain.ClassroomAcademic{
			ID: u.ClassroomID,
		},
		CountPlan:		u.CountPlan,
		MaxPoint:		u.MaxPoint,
	}

	err := h.cspUsecase.Store(c, &csp)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	u.ID = csp.ID
	c.JSON(http.StatusOK, api.ResponseSuccess("create classroom subject plan success", u))
}

func (h cspHandler) Update(c *gin.Context) {
	var u dto.ClassroomSubjectPlanRequest
	if err := c.ShouldBindJSON(&u); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(domain.ErrUnprocess.Error(), err_code),
        )
        return
    }
    currentRoleUser := c.GetInt("role")
    currentUserID := c.GetInt64("user_id")

    if currentRoleUser == domain.RoleTeacher {
        u.TeacherID = currentUserID
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

	csp := domain.ClassroomSubjectPlan {
		ID:				u.ID,
		Type: 			u.Type,
		Name:			u.Name,
		Desc:			u.Desc,
		Teacher:		domain.User{
			ID:	u.TeacherID,
		},
		Subject:		domain.ClassroomSubject{
			ID: u.SubjectID,
		},
		Classroom: 		domain.ClassroomAcademic{
			ID: u.ClassroomID,
		},
		CountPlan:		u.CountPlan,
		MaxPoint:		u.MaxPoint,
	}

	err := h.cspUsecase.Update(c, &csp)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	c.JSON(http.StatusOK, api.ResponseSuccess("update classroom subject plan success", u))
}

func (h cspHandler) Delete(c *gin.Context) {
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
	
	err = h.cspUsecase.Delete(c, int64(id))
	if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	c.JSON(http.StatusOK, api.ResponseSuccess("delete classroom subject plan success", nil))
}