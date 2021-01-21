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
	"github.com/shellrean/extraordinary-raport/interface/http/api"
	"github.com/shellrean/extraordinary-raport/interface/http/middleware"
)

type studentHandler struct {
	studentUsecase		domain.StudentUsecase
	config				*config.Config
	mddl 				*middleware.GoMiddleware
}

func NewStudentHandler(r *gin.Engine, m domain.StudentUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
	handler := &studentHandler{
		studentUsecase:	m,
		config: cfg,
		mddl:	mddl,
	}
    student := r.Group("/students")
    student.Use(handler.mddl.Auth())

    student.GET("/", handler.Index)
    student.GET("/:id", handler.Show)
	student.POST("/", handler.Store)
    student.PUT("/:id", handler.Update)
    student.DELETE("/:id", handler.Delete)
}

func (h *studentHandler) Index(c *gin.Context) {
	limS, _ := c.GetQuery("limit")
	lim, _ := strconv.Atoi(limS)
    cursor, _ := c.GetQuery("cursor")
    query, _ := c.GetQuery("q")

	res, nextCursor, err := h.studentUsecase.Fetch(c, query, cursor, int64(lim))
	if err != nil {
		error_code := helper.GetErrorCode(err)
		c.JSON(
			api.GetHttpStatusCode(err), 
			api.ResponseError(err.Error(), error_code),
		)
		return
    }
    
    var data []dto.StudentResponse
    for _, item := range res {
        student := dto.StudentResponse{
            ID:             item.ID,
            SRN:            item.SRN,
            NSRN:           item.NSRN,
            Name:           item.Name,
            Gender:         item.Gender,
            BirthPlace:     item.BirthPlace,
            BirthDate:      item.BirthDate,
            ReligionID:     item.Religion.ID,
            Address:        item.Address,
            Telp:           item.Telp,
            SchoolBefore:   item.SchoolBefore,
            AcceptedGrade:  item.AcceptedGrade,
            AcceptedDate:   item.AcceptedDate,
            FamillyStatus:  item.Familly.Status,
            FamillyOrder:   item.Familly.Order,
            FatherName:     item.Father.Name,
            FatherAddress:  item.Father.Address,
            FatherProfession: item.Father.Profession,
            FatherTelp:     item.Father.Telp,
            MotherName:     item.Mother.Name,
            MotherAddress:  item.Mother.Address,
            MotherProfession: item.Mother.Profession,
            MotherTelp:     item.Mother.Telp,
            GrdName:        item.Guardian.Name,
            GrdAddress:     item.Guardian.Address,
            GrdProfession:  item.Guardian.Profession,
            GrdTelp:        item.Guardian.Telp,
        }
        data = append(data, student)
    }

	c.Header("X-Cursor", nextCursor)
	c.JSON(http.StatusOK, api.ResponseSuccess("success",data))
}

func (h *studentHandler) Show(c *gin.Context) {
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
    res, err := h.studentUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    if res == (domain.Student{}) {
        err_code := helper.GetErrorCode(domain.ErrNotFound)
        c.JSON(
            http.StatusNotFound,
            api.ResponseError(domain.ErrNotFound.Error(), err_code),
        )
        return
    }

    data := dto.StudentResponse{
        ID:             res.ID,
        SRN:            res.SRN,
        NSRN:           res.NSRN,
        Name:           res.Name,
        Gender:         res.Gender,
        BirthPlace:     res.BirthPlace,
        BirthDate:      res.BirthDate,
        ReligionID:     res.Religion.ID,
        Address:        res.Address,
        Telp:           res.Telp,
        SchoolBefore:   res.SchoolBefore,
        AcceptedGrade:  res.AcceptedGrade,
        AcceptedDate:   res.AcceptedDate,
        FamillyStatus:  res.Familly.Status,
        FamillyOrder:   res.Familly.Order,
        FatherName:     res.Father.Name,
        FatherAddress:  res.Father.Address,
        FatherProfession: res.Father.Profession,
        FatherTelp:     res.Father.Telp,
        MotherName:     res.Mother.Name,
        MotherAddress:  res.Mother.Address,
        MotherProfession: res.Mother.Profession,
        MotherTelp:     res.Mother.Telp,
        GrdName:        res.Guardian.Name,
        GrdAddress:     res.Guardian.Address,
        GrdProfession:  res.Guardian.Profession,
        GrdTelp:        res.Guardian.Telp,
    }

    c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}

func (h *studentHandler) Store(c *gin.Context) {
	var u dto.StudentResponse
    if err := c.ShouldBindJSON(&u); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(err.Error(), err_code),
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
    
    student := domain.Student {
        ID:             u.ID,     
        SRN:            u.SRN,
        NSRN:           u.NSRN,			
        Name:           u.Name,
        Gender:         u.Gender,
        BirthPlace:     u.BirthPlace,
        BirthDate:      u.BirthDate,
        Religion:       domain.Religion{ID: u.ReligionID},
        Address:        u.Address,
        Telp:           u.Telp,
        SchoolBefore:   u.SchoolBefore,
        AcceptedGrade:  u.AcceptedGrade,
        AcceptedDate:   u.AcceptedDate,
        Familly:        domain.Familly{Status: u.FamillyStatus, Order: u.FamillyOrder},
        Father:         domain.Parent{Name: u.FatherName, Address: u.FatherAddress, Profession: u.FatherProfession, Telp: u.FatherTelp},
        Mother:         domain.Parent{Name: u.MotherName, Address: u.MotherAddress, Profession: u.MotherProfession, Telp: u.MotherTelp},			
        Guardian:       domain.Parent{Name: u.GrdName, Address: u.GrdAddress, Profession: u.GrdProfession, Telp: u.GrdTelp},
    }
	
	err := h.studentUsecase.Store(c, &student)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    u.ID = student.ID
    c.JSON(http.StatusOK, api.ResponseSuccess("create student success", u))
}

func (h *studentHandler) Update(c *gin.Context) {
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
    res, err := h.studentUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    if res == (domain.Student{}) {
        err_code := helper.GetErrorCode(domain.ErrNotFound)
        c.JSON(
            http.StatusNotFound,
            api.ResponseError(domain.ErrNotFound.Error(), err_code),
        )
        return
    }
	
	var u dto.StudentResponse
    if err := c.ShouldBindJSON(&u); err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(err.Error(), err_code),
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
    u.ID = int64(id)
    student := domain.Student {
        ID:             u.ID,     
        SRN:            u.SRN,
        NSRN:           u.NSRN,			
        Name:           u.Name,
        Gender:         u.Gender,
        BirthPlace:     u.BirthPlace,
        BirthDate:      u.BirthDate,
        Religion:       domain.Religion{ID: u.ReligionID},
        Address:        u.Address,
        Telp:           u.Telp,
        SchoolBefore:   u.SchoolBefore,
        AcceptedGrade:  u.AcceptedGrade,
        AcceptedDate:   u.AcceptedDate,
        Familly:        domain.Familly{Status: u.FamillyStatus, Order: u.FamillyOrder},
        Father:         domain.Parent{Name: u.FatherName, Address: u.FatherAddress, Profession: u.FatherProfession, Telp: u.FatherTelp},
        Mother:         domain.Parent{Name: u.MotherName, Address: u.MotherAddress, Profession: u.MotherProfession, Telp: u.MotherTelp},			
        Guardian:       domain.Parent{Name: u.GrdName, Address: u.GrdAddress, Profession: u.GrdProfession, Telp: u.GrdTelp},
    }
	err = h.studentUsecase.Update(c, &student)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("update student success", u))
}

func (h *studentHandler) Delete(c *gin.Context) {
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
    res, err := h.studentUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    if res == (domain.Student{}) {
        err_code := helper.GetErrorCode(domain.ErrNotFound)
        c.JSON(
            http.StatusNotFound,
            api.ResponseError(domain.ErrNotFound.Error(), err_code),
        )
        return
    }

    err = h.studentUsecase.Delete(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("delete student success", make([]string,0)))
}