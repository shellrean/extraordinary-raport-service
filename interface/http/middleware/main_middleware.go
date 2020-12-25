package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/shellrean/extraordinary-raport/domain"
	"github.com/shellrean/extraordinary-raport/config"
	"github.com/shellrean/extraordinary-raport/entities/helper"
	"github.com/shellrean/extraordinary-raport/interface/http/api"
)

type GoMiddleware struct {
	cfg		*config.Config
} 

func InitMiddleware(cfg *config.Config) *GoMiddleware {
	return &GoMiddleware{
		cfg:	cfg,
	}
}

func (m *GoMiddleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func (m *GoMiddleware) Auth() gin.HandlerFunc{
	return func(c *gin.Context) {
		tokenString := helper.ExtractToken(c.GetHeader("Authorization"))
		if tokenString == "" {
			c.AbortWithStatusJSON(
				api.GetHttpStatusCode(domain.ErrHeaderMiss),
				api.ResponseError(domain.ErrHeaderMiss.Error(), helper.GetErrorCode(domain.ErrHeaderMiss)),
			)
			return
		}
		token, err := helper.VerifyToken(m.cfg.JWTAccessKey, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(
				api.GetHttpStatusCode(err),
				api.ResponseError(err.Error(), helper.GetErrorCode(err)),
			)
			return
		}
		err = helper.TokenValid(token)
		if err != nil {
			c.AbortWithStatusJSON(
				api.GetHttpStatusCode(err), 
				api.ResponseError(err.Error(), helper.GetErrorCode(err)),
			)
			return
		}
		data := helper.ExtractTokenMetadata(token)
		c.Set("user-meta", data)
		c.Next()
	}
}