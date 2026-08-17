package handle

import (
	service "Stash/internal/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

func FindFiles(c *gin.Context) {
	userID := c.GetString("userID")
	value := c.Query("value")
	TypeFind := c.Query("TypeFind")
	if TypeFind == "MyDrive" {
		isFoundFiles, isFoundFolders, errCall := service.FindingHome(userID, value, TypeFind)
		if errCall != nil {
			fmt.Println(errCall)
			c.JSON(400, gin.H{
				"Error":   true,
				"Message": errCall.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Error":              false,
			"List_Found":         isFoundFiles,
			"List_Found_Folders": isFoundFolders,
		})
		return
	}
	if TypeFind == "Star" {
		isFoundFiles, isFoundFolders, errCall := service.FindingStar(userID, value, TypeFind)
		if errCall != nil {
			fmt.Println(errCall)
			c.JSON(400, gin.H{
				"Error":   true,
				"Message": errCall.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Error":              false,
			"List_Found":         isFoundFiles,
			"List_Found_Folders": isFoundFolders,
		})
		return
	}
	if TypeFind == "Trash" {
		isFoundFiles, isFoundFolders, errCall := service.FindingTrash(userID, value, TypeFind)
		if errCall != nil {
			fmt.Println(errCall)
			c.JSON(400, gin.H{
				"Error":   true,
				"Message": errCall.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Error":              false,
			"List_Found":         isFoundFiles,
			"List_Found_Folders": isFoundFolders,
		})
		return
	}
	c.Abort()
}
