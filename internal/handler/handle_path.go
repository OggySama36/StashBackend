package handle

import (
	"Stash/internal/repository"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetPathFile(c *gin.Context) {
	userID := c.GetString("userID")
	ParentFolder := c.Query("parent_id")
	result, errCall := repository.GetFolderPath(ParentFolder, userID)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error": false,
		"Path":  result,
	})
}
