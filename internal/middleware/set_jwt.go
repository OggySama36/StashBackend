package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetToken(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(
		"JWT",
		token,
		3600,
		"/",
		"",
		true,
		true,
	)
}

func ExpireToken(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(
		"JWT",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)
}
