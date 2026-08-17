package middleware

import (
	"Stash/config"
	"fmt"

	keyfunc "github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthHandler(c *gin.Context) {
	supabase_url := config.App.SupabaseURL + "/auth/v1/.well-known/jwks.json"

	tokenStrings, errGetCookie := c.Cookie("JWT")
	if errGetCookie != nil || tokenStrings == "" {
		c.JSON(401, gin.H{"Error": true})
		c.Abort()
		return
	}

	jwks, err := keyfunc.NewDefault([]string{supabase_url})
	if err != nil {
		c.JSON(401, gin.H{"Error": true})
		c.Abort()
		return
	}
	token, errParse := jwt.Parse(tokenStrings, jwks.Keyfunc)
	if errParse != nil || !token.Valid {
		fmt.Println("Token has been expired")
		c.JSON(401, gin.H{"Error": true, "Type": "Token expired"})
		c.Abort()
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	c.Set("userID", claims["sub"])
	c.Set("email", claims["email"])
	c.Next()
}

func GetMe(c *gin.Context) {
	var Result []map[string]interface{}
	_, errFind := config.SupabaseAdmin.From("profiles").
		Select("username, create_at, avatar, theme", "", false).
		Eq("email", c.GetString("email")).
		Eq("id", c.GetString("userID")).
		ExecuteTo(&Result)

	if errFind != nil {
		fmt.Println("Failed to find:", errFind)
		return
	}
	if len(Result) == 0 {
		fmt.Println("No username found!")
		return
	}

	NameUser, _ := Result[0]["username"].(string)
	JoinedDate, _ := Result[0]["create_at"].(string)
	AvatarLink, _ := Result[0]["avatar"].(string)
	Theme, _ := Result[0]["theme"].(string)
	c.JSON(200, gin.H{
		"Error":      false,
		"Email":      c.GetString("email"),
		"Name":       NameUser,
		"JoinedDate": JoinedDate,
		"AvatarLink": AvatarLink,
		"Theme":      Theme,
	})
}
