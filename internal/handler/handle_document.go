package handle

import (
	"Stash/internal/repository"
	service "Stash/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type request_document struct {
	FolderName   string `json:"Name_Folder"`
	ParentFolder string `json:"parent_folder"`
}

func CreateFolder(c *gin.Context) {
	var body request_document
	userID := c.GetString("userID")
	errParse := c.ShouldBindJSON(&body)
	if errParse != nil {
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errParse.Error(),
		})
		return
	}
	resultResponse, errCall := service.Folder(userID, body.FolderName, body.ParentFolder)
	if errCall != nil {
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error":        false,
		"Message":      resultResponse.Message,
		"Folder_id":    resultResponse.Folder_id,
		"Name_Folder":  body.FolderName,
		"ParentFolder": resultResponse.ParentFolder,
	})
}

func SaveFile(c *gin.Context) {
	userID := c.GetString("userID")
	file, errGetFile := c.FormFile("file")
	if errGetFile != nil {
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errGetFile.Error(),
		})
		return
	}
	parent_folder := c.PostForm("parent_folder")
	resultSave, errCall := service.File(file, parent_folder, userID)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error":        false,
		"File_id":      resultSave.File_id,
		"Name":         file.Filename,
		"ParentFolder": parent_folder,
		"MimeType":     resultSave.MimeType,
		"Size":         resultSave.Size,
	})
}

func Load_Folders(c *gin.Context) {
	userID := c.GetString("userID")
	result_load_folder, errCall := repository.Folders_Loading(userID)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error":        false,
		"Type":         "Folders",
		"Message":      "Load success",
		"List_folders": result_load_folder,
	})
}

func Load_Files(c *gin.Context) {
	userID := c.GetString("userID")
	result_load_file, errCall := repository.Files_Loading(userID)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error":      false,
		"Type":       "Files",
		"Message":    "Load success",
		"List_files": result_load_file,
	})
}

func Load_Trashes_Files(c *gin.Context) {
	userID := c.GetString("userID")
	result_load_trashes, errCall := repository.Trashes_Files_Loading(userID)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
	}
	c.JSON(200, gin.H{
		"Error":      false,
		"Type":       "Trashes",
		"Message":    "Load success",
		"List_files": result_load_trashes,
	})
}

func Load_Stars_Files(c *gin.Context) {
	userID := c.GetString("userID")
	result_load_file, errCall := repository.Stars_Files_Loading(userID)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error":      false,
		"Type":       "Star",
		"Message":    "Load success",
		"List_files": result_load_file,
	})
}

func Load_Stars_Folders(c *gin.Context) {
	userID := c.GetString("userID")
	result_load_stars, errCall := repository.Stars_Folders_Loading(userID)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error":        false,
		"Type":         "Star",
		"Message":      "Load success",
		"List_folders": result_load_stars,
	})
}

func Load_Trashes_Folders(c *gin.Context) {
	userID := c.GetString("userID")
	result_load_trashes, errCall := repository.Trashes_Folders_Loading(userID)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
	}
	c.JSON(200, gin.H{
		"Error":        false,
		"Type":         "Trashes",
		"Message":      "Load success",
		"List_folders": result_load_trashes,
	})
}

func Create_URL(c *gin.Context) {
	userID := c.GetString("userID")
	File_id := c.Query("file_id")
	result_url, errCall := service.Get_Url(userID, File_id)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error":   false,
		"File_id": result_url.File_id,
		"URL":     result_url.URL,
	})
}

func RemoveFile(c *gin.Context) {
	userID := c.GetString("userID")
	file_id := c.Query("file_id")
	Type_handle := c.Query("type")
	if Type_handle == "Remove" {
		Result, errCall := service.RemoveFile(userID, file_id)
		if errCall != nil {
			fmt.Println(errCall)
			c.JSON(400, gin.H{
				"Error":   true,
				"Message": errCall.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Error":   false,
			"Message": Result.Message,
		})
		c.Abort()
	}
	if Type_handle == "Restore" {
		Result, errCall := service.RestoreFile(userID, file_id)
		if errCall != nil {
			fmt.Println(errCall)
			c.JSON(400, gin.H{
				"Error":   true,
				"Message": errCall.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Error":   false,
			"Message": Result.Message,
		})
	}
	c.Abort()
}

func RemoveFolder(c *gin.Context) {
	userID := c.GetString("userID")
	folder_id := c.Query("folder_id")
	Type_handle := c.Query("type")
	if Type_handle == "Remove" {
		Result, errCall := service.RemoveFolder(userID, folder_id)
		if errCall != nil {
			fmt.Println(errCall)
			c.JSON(400, gin.H{
				"Error":   true,
				"Message": errCall.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Error":     false,
			"Message":   Result.Message,
			"TotalSize": Result.TotalSize,
			"Folders":   Result.Folders,
			"Files":     Result.Files,
		})
		c.Abort()
	}
	if Type_handle == "Restore" {
		Result, errCall := service.RestoreFolder(userID, folder_id)
		if errCall != nil {
			fmt.Println(errCall)
			c.JSON(400, gin.H{
				"Error":   true,
				"Message": errCall.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Error":     false,
			"Message":   Result.Message,
			"TotalSize": Result.TotalSize,
			"Folders":   Result.Folders,
			"Files":     Result.Files,
		})
	}
	c.Abort()
}

func Star_File(c *gin.Context) {
	userID := c.GetString("userID")
	file_id := c.Query("file_id")
	Type_handle := c.Query("type")
	if Type_handle == "Star" {
		Result_Star, errCall := service.StarFile(userID, file_id)
		if errCall != nil {
			fmt.Println(errCall)
			c.JSON(400, gin.H{
				"Error":   true,
				"Message": errCall.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Error":   false,
			"Message": Result_Star.Message,
		})
		c.Abort()
	}
	if Type_handle == "Unstar" {
		Result_Star, errCall := service.UnstarFile(userID, file_id)
		if errCall != nil {
			fmt.Println(errCall)
			c.JSON(400, gin.H{
				"Error":   true,
				"Message": errCall.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Error":   false,
			"Message": Result_Star.Message,
		})
	}
	c.Abort()
}

func Star_Folder(c *gin.Context) {
	userID := c.GetString("userID")
	folder_id := c.Query("folder_id")
	Type_handle := c.Query("type")
	if Type_handle == "Star" {
		Result_Star, errCall := service.StarFolder(userID, folder_id)
		if errCall != nil {
			fmt.Println(errCall)
			c.JSON(400, gin.H{
				"Error":   true,
				"Message": errCall.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Error":   false,
			"Message": Result_Star.Message,
		})
		c.Abort()
	}
	if Type_handle == "Unstar" {
		Result_Star, errCall := service.UnstarFolder(userID, folder_id)
		if errCall != nil {
			fmt.Println(errCall)
			c.JSON(400, gin.H{
				"Error":   true,
				"Message": errCall.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Error":   false,
			"Message": Result_Star.Message,
		})
	}
	c.Abort()
}

func Delete_Permanently_Files(c *gin.Context) {
	userID := c.GetString("userID")
	file_id := c.Query("file_id")
	result, errCall := service.Delete_Files(userID, file_id)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error":   false,
		"Message": result,
	})
}

func Delete_Permanently_Folders(c *gin.Context) {
	userID := c.GetString("userID")
	folder_id := c.Query("folder_id")
	errCall := service.Delete_Folders(userID, folder_id)
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
	})
}

type RequestShare struct {
	File_id  string `json:"file_id"`
	FileName string `json:"fileName"`
	Note     string `json:"note"`
}

func ShareFileURL(c *gin.Context) {
	sender := c.GetString("email")
	userID := c.GetString("userID")
	var dataShare RequestShare
	errParse := c.ShouldBindJSON(&dataShare)
	if errParse != nil {
		fmt.Println(errParse)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errParse.Error(),
		})
		return
	}
	Result, errCall := service.Copy_url_share(dataShare.File_id, sender, userID)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error":    false,
		"Share_id": Result.Share_id,
		"File_id":  Result.File_id,
		"Name":     Result.Name,
		"MimeType": Result.MimeType,
		"SharedBy": Result.SharedBy,
		"SharedTo": Result.SharedTo,
		"Method":   Result.Method,
	})
}

func ShareFileGmail(c *gin.Context) {
	var dataShare RequestShare
	userID := c.GetString("userID")
	errParse := c.ShouldBindJSON(&dataShare)
	if errParse != nil {
		fmt.Println(errParse)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errParse.Error(),
		})
		return
	}
	sender := c.GetString("email")
	recipient := c.Query("recipient")
	fmt.Println(recipient)
	Result, errCall := service.Send_gmail_file(dataShare.File_id, dataShare.FileName, dataShare.Note, sender, recipient, userID)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error":    false,
		"Share_id": Result.Share_id,
		"File_id":  Result.File_id,
		"Name":     Result.Name,
		"MimeType": Result.MimeType,
		"SharedBy": Result.SharedBy,
		"SharedTo": Result.SharedTo,
		"Method":   Result.Method,
	})
}

type RenameReq struct {
	NewName string `json:"new_name"`
}

func RenameThisFolder(c *gin.Context) {
	var Newname RenameReq
	userID := c.GetString("userID")
	folderID := c.Query("folder_id")
	errAssign := c.ShouldBindJSON(&Newname)
	fmt.Println(Newname.NewName)
	if errAssign != nil {
		fmt.Println(errAssign)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errAssign.Error(),
		})
		return
	}
	result, errCall := service.Rename_folder(userID, folderID, Newname.NewName)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error":   false,
		"Message": result.Message,
	})
}

func LoadSharedFiles(c *gin.Context) {
	userID := c.GetString("userID")
	email := c.GetString("email")
	ListShared, errCall := service.LoadShared(userID, email)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error":      false,
		"ListShared": ListShared,
	})
}
