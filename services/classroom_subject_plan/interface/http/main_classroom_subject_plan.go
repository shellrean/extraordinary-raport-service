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
	cspUsecase	domain.ClassroomSubjectPlanUsecase
	config 		*config.Config
	mddl 		*middleware.GoMiddleware
}

func New(r *gin.Engine, m domain.ClassroomSubjectPlanUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	h := &handler {
		cspUsecase:	m,
		config:		cfg,
		mddl:		mddl,
	}
	csp := r.Group("/classroom-subject-plans")
	csp.Use(h.mddl.Auth())

    csp.POST("/", h.mddl.MustRole([]int{domain.RoleTeacher, domain.RoleAdmin}), h.Fetch)
	csp.POST("/classroom-subject-plan", h.mddl.MustRole([]int{domain.RoleTeacher}), h.Store)
	csp.PUT("/classroom-subject-plan", h.mddl.MustRole([]int{domain.RoleTeacher}), h.Update)
    csp.DELETE("/classroom-subject-plan/:id", h.mddl.MustRole([]int{domain.RoleTeacher}), h.Delete)
    csp.DELETE("delete", h.mddl.MustRole([]int{domain.RoleTeacher}), h.DeleteMultiple)
    csp.GET("/classroom-subject-plan/:id", h.mddl.MustRole([]int{domain.RoleTeacher, domain.RoleAdmin}), h.Show)
}

func (h *handler) Fetch(c *gin.Context) {
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

func (h *handler) Show(c *gin.Context) {
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

    res, err := h.cspUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    currentRoleUser := c.GetInt("role")
    currentUserID := c.GetInt64("user_id")

    if currentRoleUser == domain.RoleTeacher && currentUserID != res.Teacher.ID {
        c.JSON(
			api.GetHttpStatusCode(domain.ErrNoAuthorized),
			api.ResponseError(domain.ErrNoAuthorized.Error(), helper.GetErrorCode(domain.ErrNoAuthorized)),
        )
        return
    }

    data := dto.ClassroomSubjectPlanResponse {
        ID: res.ID,
        Type: res.Type,
        Name: res.Name,
        Desc: res.Desc,
        TeacherID: res.Teacher.ID,
        SubjectID: res.Subject.ID,
        SubjectName: res.Subject.Subject.Name,
        ClassroomID: res.Classroom.ID,
        CountPlan: res.CountPlan,
        MaxPoint: res.MaxPoint,
    }

    c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}

func (h *handler) Store(c *gin.Context) {
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

func (h *handler) Update(c *gin.Context) {
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

func (h *handler) Delete(c *gin.Context) {
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

func (h *handler) DeleteMultiple(c *gin.Context) {
    query, _ := c.GetQuery("q")
    if query == "" {
        err_code := helper.GetErrorCode(domain.ErrBadParamInput)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseError(domain.ErrBadParamInput.Error(), err_code),
        )
        return
    }

    err := h.cspUsecase.DeleteMultiple(c, query)
	if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	c.JSON(http.StatusOK, api.ResponseSuccess("delete multiple classroom subject plan success", nil))
}