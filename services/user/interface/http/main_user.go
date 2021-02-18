package handler

import (
    "net/http"
    "strings"
    "strconv"
    "os"
    "io"
    "path/filepath"

    "github.com/google/uuid"
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
    userUsecase     domain.UserUsecase
    config          *config.Config
    mddl            *middleware.GoMiddleware
}

func New(r *gin.Engine, m domain.UserUsecase, cfg *config.Config, mddl *middleware.GoMiddleware) {
    h := &handler{
        userUsecase:    m,
        config:         cfg,
        mddl:           mddl,
    }

    user := r.Group("/users")
    user.Use(h.mddl.Auth())

    user.GET("/", h.FetchUsers)
    user.GET("user/:id", h.Show)
    user.POST("user", h.Store)
    user.PUT("user/:id", h.Update)
    user.DELETE("user/:id", h.Delete)
    user.DELETE("delete", h.DeleteMultiple)
    user.POST("import", h.Import)
    user.GET("auth-current", h.CurrentLogin)

    r.POST("/auth", h.Autheticate)
    r.POST("/refresh-token", h.RefreshToken)
}

func (h *handler) FetchUsers(c *gin.Context) {
    limS , _ := c.GetQuery("limit")
    lim, _ := strconv.Atoi(limS)
    cursor, _ := c.GetQuery("cursor")
    query, _ := c.GetQuery("q")

    res, nextCursor, err := h.userUsecase.Fetch(c, query, cursor, int64(lim))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    var data []dto.UserResponse
    for _, item := range res {
        user := dto.UserResponse{
            ID:     item.ID,
            Name:   item.Name,
            Email:  item.Email,
            Role:   item.Role,
        }
        data = append(data, user)
    }

    c.Header("X-Cursor", nextCursor)
    c.JSON(http.StatusOK, api.ResponseSuccess("success", data)) 
}

func (h *handler) Autheticate(c *gin.Context) {
    var u dto.UserLogin
    if err := c.ShouldBindJSON(&u); err != nil {
        error_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(domain.ErrUnprocess.Error(), error_code),
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
        error_code := helper.GetErrorCode(domain.ErrValidation)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseErrorWithData(domain.ErrValidation.Error(), error_code, reserr),
        )
        return
    }

    user := domain.User{
        Email: u.Email,
        Password: u.Password,
    }

    res, err := h.userUsecase.Authentication(c, user)
    if err != nil {
        error_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), error_code),
        )
        return
    }

    data := dto.TokenResponse {
        AccessToken: res.AccessToken,
        RefreshToken: res.RefreshToken,
    }
    
    c.JSON(http.StatusOK, api.ResponseSuccess("success",data))
}

func (h *handler) RefreshToken(c *gin.Context) {
    var u dto.TokenResponse
    if err := c.ShouldBindJSON(&u); err != nil {
        error_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity, 
            api.ResponseError(domain.ErrUnprocess.Error(), error_code),
        )
        return
    }

    token := domain.Token {
        AccessToken: u.AccessToken,
        RefreshToken: u.RefreshToken,
    }

    err := h.userUsecase.RefreshToken(c, &token)
    if err != nil {
        error_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err), 
            api.ResponseError(err.Error(), error_code),
        )
        return
    }

    data := dto.TokenResponse{
        token.AccessToken,
        token.RefreshToken,
    }
    
    c.JSON(http.StatusOK, api.ResponseSuccess("success",data))
}

func (h *handler) CurrentLogin(c *gin.Context) {
    userID := c.GetInt64("user_id")

    res, err := h.userUsecase.UserCurrentLogin(c, userID)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    
    data := dto.UserResponse {
        ID: res.ID,
        Name: res.Name,
        Email: res.Email,
        Role: res.Role,
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
    res := domain.User{}
    res, err = h.userUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    data := dto.UserResponse {
        ID: res.ID,
        Name: res.Name,
        Email: res.Email,
        Role: res.Role,
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("success", data))
}

func (h *handler) Store(c *gin.Context) {
    var u dto.UserCreatePayload
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

    user := domain.User{
        Name:           u.Name,
        Email:          u.Email,
        Role:           u.Role,
        Password:       u.Password,
    }

    err := h.userUsecase.Store(c, &user)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    data := dto.UserResponse{
        ID:     user.ID,
        Name:   user.Name,
        Email:  user.Email,
        Role:   user.Role,
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("create user success", data))
}

func (h *handler) Import(c *gin.Context) {
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        err_code := helper.GetErrorCode(domain.ErrUnprocess)
        c.JSON(
            http.StatusUnprocessableEntity,
            api.ResponseError(domain.ErrUnprocess.Error(), err_code),
        )
        return
    }

    filename := header.Filename

    fileExtension := filepath.Ext(filename)
    if (fileExtension != ".xls" && fileExtension != ".xlsx") {
        err_code := helper.GetErrorCode(domain.ErrFileNotAllowed)
        c.JSON(
            http.StatusBadRequest,
            api.ResponseError(domain.ErrFileNotAllowed.Error(), err_code),
        )
        return
    }

    fullPathFile := filepath.Join("storage", "app", "_tmp", uuid.NewString()+filename)

    out, err := os.Create(fullPathFile)
    if err != nil {
        err_code := helper.GetErrorCode(domain.ErrServerError)
        c.JSON(
            http.StatusInternalServerError,
            api.ResponseError(domain.ErrServerError.Error(), err_code),
        )
        return
    }
    defer os.Remove(fullPathFile)
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        err_code := helper.GetErrorCode(domain.ErrServerError)
        c.JSON(
            http.StatusInternalServerError,
            api.ResponseError(domain.ErrServerError.Error(), err_code),
        )
        return
    }

    err = h.userUsecase.ImportFromExcel(c, fullPathFile)
    
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    c.JSON(http.StatusOK, api.ResponseSuccess("import user success", nil))
}

func (h *handler) Update(c *gin.Context) {
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
    res, err := h.userUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    if res == (domain.User{}) {
        err_code := helper.GetErrorCode(domain.ErrNotFound)
        c.JSON(
            api.GetHttpStatusCode(domain.ErrNotFound),
            api.ResponseError(domain.ErrNotFound.Error(), err_code),
        )
        return
    }

    var u dto.UserResponse
    err = c.ShouldBindJSON(&u)
    if err != nil {
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
    u.ID = int64(id)
    user := domain.User{
        ID:     u.ID,
        Name:   u.Name,
        Email:  u.Email,
        Role:   u.Role,
    }
    err = h.userUsecase.Update(c, &user)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("update user success", u))
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
    res, err := h.userUsecase.GetByID(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }

    if res == (domain.User{}) {
        err_code := helper.GetErrorCode(domain.ErrNotFound)
        c.JSON(
            api.GetHttpStatusCode(domain.ErrNotFound),
            api.ResponseError(domain.ErrNotFound.Error(), err_code),
        )
        return
    }

    err = h.userUsecase.Delete(c, int64(id))
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("delete user success", make([]string,0)))
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
    
    err := h.userUsecase.DeleteMultiple(c, query)
    if err != nil {
        err_code := helper.GetErrorCode(err)
        c.JSON(
            api.GetHttpStatusCode(err),
            api.ResponseError(err.Error(), err_code),
        )
        return
    }
    c.JSON(http.StatusOK, api.ResponseSuccess("delete multiple user success", make([]string,0)))
}