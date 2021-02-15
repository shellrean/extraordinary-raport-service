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

type sprUsecase domain.ClassroomSubjectPlanResultUsecase
type sprResponse dto.ClassroomSubjectPlanResultResponse
type sprRequest dto.ClassroomSubjectPlanResultRequest

type sprHandler struct {
	sprUsecase 	sprUsecase
	cfg 		*config.Config
	mddl 		*middleware.GoMiddleware
}

func NewClassroomSubjectPlanResultHandler(r *gin.Engine, u sprUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	handler := &sprHandler{
		sprUsecase:	u,
		cfg:		cfg,
		mddl:		mddl,
	}
	spr := r.Group("/cspr")
	spr.Use(handler.mddl.Auth())

	spr.POST("s", handler.Store) // Store single plan result
	spr.GET("plan/:id", handler.FetchByPlan)
}

func (h sprHandler) FetchByPlan(c *gin.Context) {
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

	res, err := h.sprUsecase.FetchByPlan(c, int64(id))
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}

	var data []dto.ClassroomSubjectPlanResultResponse
	for _, item := range res {
		spr := dto.ClassroomSubjectPlanResultResponse{
			ID:			item.ID,
			Index:		item.Index,
			StudentID:	item.Student.ID,
			SubjectID:	item.Subject.ID,
			PlanID:		item.Plan.ID,
			Number:	 	item.Number,
			UpdatedByID:item.UpdatedBy.ID,
		}
		data = append(data, spr)
	}

	c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}	

func (h sprHandler) Store(c *gin.Context) {
	var u sprRequest
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

	sess, _ := c.Get("user_id")
	user_id := sess.(int64)

	s := domain.ClassroomSubjectPlanResult {
		Index: u.Index,
		Student: domain.ClassroomStudent{
			ID: u.StudentID,
		},
		Subject: domain.ClassroomSubject{
			ID: u.SubjectID,
		},
		Plan: domain.ClassroomSubjectPlan{
			ID: u.PlanID,
		},
		Number: u.Number,
		UpdatedBy: domain.User{
			ID: user_id,
		},
	}

	err := h.sprUsecase.Store(c, &s)
	if err != nil {
		err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
	}
	u.ID = s.ID
	c.JSON(http.StatusOK, api.ResponseSuccess("create plan result success", u))
}

