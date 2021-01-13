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

type subjectHandler struct {
    subjectUsecase  domain.SubjectUsecase
    config          *config.Config
    mddl            *middleware.GoMiddleware
}

func NewSubjectHandler(r *gin.Engine, m domain.SubjectUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	handler := &subjectHandler{
		subjectUsecase: 	m,
		config:				cfg,
		mddl:				mddl,
	}

	subject := r.Group("/subjects")
	subject.Use(handler.mddl.Auth())

    subject.GET("/", handler.Fetch)
    subject.GET("/:id", handler.Show)
}

func (h *subjectHandler) Fetch(c *gin.Context) {
	limS , _ := c.GetQuery("limit")
    lim, _ := strconv.Atoi(limS)
    cursor, _ := c.GetQuery("cursor")

    res, nextCursor, err := h.subjectUsecase.Fetch(c, cursor, int64(lim))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    var data []dto.SubjectResponse
    for _, item := range res {
        user := dto.SubjectResponse{
            ID:     item.ID,
			Name:   item.Name,
			Type:	item.Type,
        }
        data = append(data, user)
    }

    c.Header("X-Cursor", nextCursor)
    c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
}

func (h *subjectHandler) Show(c *gin.Context) {
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
    res := domain.Subject{}
    res, err = h.subjectUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    data := dto.SubjectResponse {
        ID: res.ID,
        Name: res.Name,
        Type: res.Type,
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}