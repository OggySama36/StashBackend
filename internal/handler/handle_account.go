package handle

import (
	"Stash/config"
	"Stash/internal/middleware"
	service "Stash/internal/services"
	"context"
	"fmt"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
	storage_go "github.com/supabase-community/storage-go"
	"gopkg.in/gomail.v2"
)

type request_account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Message  string
	UserID   string
}

func RegisterHandler(c *gin.Context) {
	var body request_account
	errParse := c.ShouldBindJSON(&body)
	if errParse != nil {
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": "Invalid Request",
		})
		return
	}
	resultResponse, errCall := service.Register(body.Username, body.Password, body.Email)
	if errCall != nil {
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	middleware.SetToken(c, resultResponse.Token)
	c.JSON(200, gin.H{
		"Error":      false,
		"Message":    "Register Successfully",
		"Email":      resultResponse.Email,
		"Name":       resultResponse.Name,
		"JoinedDate": resultResponse.JoinedDate,
		"AvatarLink": resultResponse.AvatarLink,
		"Theme":      resultResponse.Theme,
		"Token":      resultResponse.Token,
	})
}

func LoginHandler(c *gin.Context) {
	var body request_account
	errParse := c.ShouldBindJSON(&body)
	if errParse != nil {
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": "Invalid Request",
		})
		return
	}
	errCall, resultLogin := service.Login(body.Email, body.Password)
	if errCall != nil {
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}

	middleware.SetToken(c, resultLogin.Token)
	c.JSON(200, gin.H{
		"Error":      false,
		"Message":    "Login Successfully!",
		"Email":      resultLogin.Email,
		"Name":       resultLogin.Name,
		"JoinedDate": resultLogin.JoinedDate,
		"AvatarLink": resultLogin.AvatarLink,
		"Theme":      resultLogin.Theme,
		"Token":      resultLogin.Token,
	})
}

func LogoutHandler(c *gin.Context) {
	middleware.ExpireToken(c)
	c.JSON(200, gin.H{
		"Error":   false,
		"Message": "Log out succeed",
	})
}

type changePwd struct {
	OldPwd string `json:"oldPwd"`
	NewPwd string `json:"newPwd"`
}

func ChangePassword(c *gin.Context) {
	var changePwd changePwd
	userId := c.GetString("userID")
	email := c.GetString("email")
	errParse := c.ShouldBindJSON(&changePwd)
	if errParse != nil {
		fmt.Println(errParse)
		c.JSON(401, gin.H{
			"Error":   true,
			"Message": errParse.Error(),
		})
		return
	}
	if changePwd.OldPwd == changePwd.NewPwd {
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": "New password must be different from current password!",
		})
		return
	}
	errCall := service.ChangePwd(userId, email, changePwd.OldPwd, changePwd.NewPwd)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(402, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error": false,
	})
}

func UploadAvatar(c *gin.Context) {
	userID := c.GetString("userID")
	avatar, errGetAvatar := c.FormFile("file")
	if errGetAvatar != nil {
		fmt.Println(errGetAvatar)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errGetAvatar.Error(),
		})
		return
	}
	fileName := avatar.Filename
	idAvatar := fmt.Sprintf("%s-%s", userID, fileName)
	openFile, errOpen := avatar.Open()
	if errOpen != nil {
		fmt.Println(errOpen)
		c.JSON(500, gin.H{
			"Error":   true,
			"Message": errOpen.Error(),
		})
		return
	}
	defer openFile.Close()
	Cloud, errCreate := cloudinary.NewFromParams(
		config.App.CloudinaryName,
		config.App.CloudinaryAPIKey,
		config.App.CloudinarySecretKey,
	)
	if errCreate != nil {
		fmt.Println(errCreate)
		c.JSON(500, gin.H{
			"Error":   true,
			"Message": errCreate.Error(),
		})
		return
	}
	uploadResult, errUpload := Cloud.Upload.Upload(context.Background(), openFile, uploader.UploadParams{
		Folder:   "avatar",
		PublicID: idAvatar,
	})
	if uploadResult.Error.Message != "" {
		c.JSON(500, gin.H{"Error": true, "Message": uploadResult.Error.Message})
		return
	}
	if errUpload != nil {
		fmt.Println(errUpload)
		c.JSON(500, gin.H{
			"Error":   true,
			"Message": errUpload.Error(),
		})
		return
	}

	ResultSaveAvatar, errCall := service.SaveUrlAvatar(userID, uploadResult.SecureURL)
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
		"AvatarLink": ResultSaveAvatar.AvatarLink,
	})
	c.Abort()
}

func ChangeTheme(c *gin.Context) {
	userID := c.GetString("userID")
	theme := c.Query("theme")
	errCall := service.AnotherTheme(userID, theme)
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
	c.Abort()
}

func DeleteAccount(c *gin.Context) {
	userIDstring := c.GetString("userID")
	userIDuuid, errParse := uuid.Parse(userIDstring)
	if errParse != nil {
		fmt.Println(errParse)
		c.JSON(500, gin.H{
			"Error":   true,
			"Message": errParse.Error(),
		})
		return
	}
	errDeleteUser := config.SupabaseAdmin.Auth.AdminDeleteUser(types.AdminDeleteUserRequest{
		UserID: userIDuuid,
	})
	if errDeleteUser != nil {
		fmt.Println(errDeleteUser)
		c.JSON(500, gin.H{
			"Error":   true,
			"Message": errDeleteUser.Error(),
		})
		return
	}
	offset := 0
	for {
		listFilesBucket, errGetListsFile := config.SupabaseAdmin.Storage.ListFiles("Stash", userIDstring, storage_go.FileSearchOptions{
			Limit:  1000,
			Offset: offset,
		})
		if errGetListsFile != nil {
			fmt.Println(errGetListsFile)
			c.JSON(500, gin.H{
				"Error":   true,
				"Message": errGetListsFile.Error(),
			})
			return
		}
		if len(listFilesBucket) == 0 {
			break
		}

		paths := make([]string, 0, len(listFilesBucket))
		for _, value := range listFilesBucket {
			paths = append(paths, userIDstring+"/"+value.Name)
		}
		_, errRemoveListFiles := config.SupabaseAdmin.Storage.RemoveFile("Stash", paths)
		if errRemoveListFiles != nil {
			fmt.Println(errRemoveListFiles)
			c.JSON(500, gin.H{
				"Error":   true,
				"Message": errRemoveListFiles.Error(),
			})
			return
		}
		//Đỡ tốn thêm 1 lần gọi function
		if len(listFilesBucket) < 1000 {
			break
		}
	}
	errResultDeleteDB_profile := service.DeleteProfile(userIDstring)
	errResultDeleteDB_files := service.DeleteAllFiles(userIDstring)
	errResultDeleteDB_folders := service.DeleteAllFolders(userIDstring)
	if errResultDeleteDB_profile != nil {
		fmt.Println(errResultDeleteDB_profile)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errResultDeleteDB_profile.Error(),
		})
		return
	}
	if errResultDeleteDB_files != nil {
		fmt.Println(errResultDeleteDB_files)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errResultDeleteDB_files.Error(),
		})
		return
	}
	if errResultDeleteDB_folders != nil {
		fmt.Println(errResultDeleteDB_folders)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errResultDeleteDB_folders.Error(),
		})
		return
	}
	middleware.ExpireToken(c)
	c.JSON(200, gin.H{
		"Error": false,
	})
}

type NewEmail struct {
	NewEmail string `json:"email"`
}

func SendOTPemail(toEmail string, OTPcode int) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", config.App.HostEmail)
	mail.SetHeader("To", toEmail)
	mail.SetHeader("Subject", "OTP code to change email")
	mail.SetBody("text/html", fmt.Sprintf("Your OTP code is: %d", OTPcode))
	send := gomail.NewDialer(config.App.SmtpHost, 587, config.App.HostEmail, config.App.ApplicationPwd)
	return send.DialAndSend(mail)
}

func ChangeEmail(c *gin.Context) {
	userID := c.GetString("userID")
	emailUser := c.GetString("email")
	var newEmail NewEmail
	errParse := c.ShouldBindJSON(&newEmail)
	if errParse != nil {
		fmt.Println(errParse)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errParse.Error(),
		})
		return
	}
	OTP, errCall := service.SaveTemporaryEmail(userID, newEmail.NewEmail)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	SendOTPemail(emailUser, OTP.OTP)
	c.JSON(200, gin.H{
		"Error": false,
	})
}

type OTPcode struct {
	Code string `json:"vrfCode"`
}

func VerifyOTPCode(c *gin.Context) {
	var Otpcode OTPcode
	errParse := c.ShouldBindJSON(&Otpcode)
	if errParse != nil {
		fmt.Println(errParse)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errParse.Error(),
		})
		return
	}
	userID := c.GetString("userID")
	OtpcodeParsed, _ := strconv.Atoi(Otpcode.Code)
	ResultVerify, errCall := service.CheckVerifyCodeAndDecide(OtpcodeParsed, userID)
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
		"NewEmail": ResultVerify.NewEmail,
	})
}

type GetForgotPassword struct {
	Email string `json:"email"`
}
type GetOtpForgotPassword struct {
	Email string `json:"email"`
	Otp   string `json:"Otp"`
}
type GetNewPassword struct {
	Email       string `json:"email"`
	NewPassword string `json:"NewPassword"`
}

func FindForgotPassword(c *gin.Context) {
	var GetForgotEmail GetForgotPassword
	errParse := c.ShouldBindJSON(&GetForgotEmail)
	if errParse != nil {
		fmt.Println(errParse)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errParse.Error(),
		})
		return
	}
	OTP, errCall := service.FindForgotEmail(GetForgotEmail.Email)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(402, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	SendOTPemail(GetForgotEmail.Email, OTP.OTP)
	c.JSON(200, gin.H{
		"Error": false,
	})
}

func VerifyForgotPasswordCode(c *gin.Context) {
	var GetOtpForgotPassword GetOtpForgotPassword
	errParse := c.ShouldBindJSON(&GetOtpForgotPassword)
	if errParse != nil {
		fmt.Println(errParse)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errParse.Error(),
		})
		return
	}
	errCall := service.CheckCodeForgotPassword(GetOtpForgotPassword.Otp, GetOtpForgotPassword.Email)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(404, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error": false,
	})
}

func ResetPassword(c *gin.Context) {
	var GetNewPassword GetNewPassword
	errParse := c.ShouldBindJSON(&GetNewPassword)
	if errParse != nil {
		fmt.Println(errParse)
		c.JSON(400, gin.H{
			"Error":   true,
			"Message": errParse.Error(),
		})
		return
	}
	errCall := service.ResetForgotPassword(GetNewPassword.Email, GetNewPassword.NewPassword)
	if errCall != nil {
		fmt.Println(errCall)
		c.JSON(402, gin.H{
			"Error":   true,
			"Message": errCall.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Error": false,
	})
	c.Abort()
}
