package config

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func SecurityMiddleware() gin.HandlerFunc {
	sec := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		SSLRedirect:           false,
		ContentSecurityPolicy: "default-src 'self'; connect-src 'self' http://localhost:3000 http://localhost:8000",
	})

	return func(c *gin.Context) {
		err := sec.Process(c.Writer, c.Request)
		if err != nil {
			c.Abort()
			return
		}
		c.Next()
	}
}
