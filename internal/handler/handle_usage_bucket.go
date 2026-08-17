package handle

import (
	"Stash/internal/repository"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Load_Usage(c *gin.Context) {
	userID := c.GetString("userID")
	result, errCall := repository.Get_Usage(userID)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Message": errCall.Error(),
			"Error":   true,
		})
		return
	}
	c.JSON(200, gin.H{
		"Error": false,
		"Used":  result.Used,
		"Left":  result.Left,
	})
}
