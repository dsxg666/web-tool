package middleware

import (
	"github.com/dsxg666/web-tool/global"
	"github.com/dsxg666/web-tool/pkg/jwt"
	"github.com/dsxg666/web-tool/pkg/result"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := c.GetHeader("JwtToken")

		if jwtToken == "" {
			c.JSON(http.StatusUnauthorized, result.NoJwtToken)
			c.Abort()
			return
		}

		claims, ok, err := jwt.ParseJwtToken(jwtToken)
		if err != nil {
			global.Logger.Error("parse err", err)
			c.JSON(http.StatusUnauthorized, result.InvalidJwtToken)
			c.Abort()
			return
		}
		if ok {
			c.Set("userId", claims.UserId)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, result.InvalidJwtToken)
			c.Abort()
			return
		}
	}
}

func LimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		global.Logger.Infof("%s try to access website", c.Request.Host)
		if c.Request.Host != "http://121.196.245.107" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
