package ginjwt

import (
	"time"

	"l2-golang-auth/config"
	"l2-golang-auth/src/noroute"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func MwInitializer() *jwt.GinJWTMiddleware {

	//============================================================================
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte(config.JWTSignaturePrivateKey),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,

		Unauthorized: unauthorizedFunc,

		TokenLookup: "header:Authorization",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}
	// Initial middleware default setting.
	authMiddleware.MiddlewareInit()
	//============================================================================

	return authMiddleware
}

func unauthorizedFunc(c *gin.Context, code int, msg string) {
	//use 404 not found instead of 401 Unauthorized.
	noroute.NotFoundHandler(c)
}
