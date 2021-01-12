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