package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/shellrean/extraordinary-raport/config"
	"github.com/shellrean/extraordinary-raport/entities/helper"
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
		token, err := helper.VerifyToken(m.cfg.JWTAccessKey, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(helper.GetStatusCode(err), gin.H{"message": err.Error()})
			return
		}
		err = helper.TokenValid(token)
		if err != nil {
			c.AbortWithStatusJSON(helper.GetStatusCode(err), gin.H{"message": err.Error()})
			return
		}
		data := helper.ExtractTokenMetadata(token)
		c.Set("user-meta", data)
		c.Next()
	}
}