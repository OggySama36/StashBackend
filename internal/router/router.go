package router

import (
	handle "Stash/internal/handler"
	"Stash/internal/middleware"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	public := engine.Group("/api/v1")
	{
		public.POST("/auth/register", handle.RegisterHandler)
		public.POST("/auth/login", handle.LoginHandler)
		public.POST("/auth/logout", handle.LogoutHandler)
		public.POST("/auth/checkAuth", middleware.AuthHandler, middleware.GetMe)
		public.POST("/auth/change/password/forgot", handle.FindForgotPassword)
		public.POST("/verify/code/password/forgot", handle.VerifyForgotPasswordCode)
		public.POST("/auth/reset/password/forgot", handle.ResetPassword)
	}
	private := engine.Group("/api/v1")
	private.Use(middleware.AuthHandler)
	{
		private.POST("/auth/change/email", handle.ChangeEmail)
		private.POST("/verify/code/email", handle.VerifyOTPCode)
		private.POST("/auth/change/password", handle.ChangePassword)
		private.POST("/upload/folders", handle.CreateFolder)
		private.POST("/upload/files", handle.SaveFile)
		private.GET("/load/folders", handle.Load_Folders)
		private.GET("/load/files", handle.Load_Files)
		private.GET("/get/usage", handle.Load_Usage)
		private.GET("/get/url_file", handle.Create_URL)
		private.POST("/handle/remove/file", handle.RemoveFile)
		private.POST("/handle/remove/folder", handle.RemoveFolder)
		private.POST("/handle/star/file", handle.Star_File)
		private.POST("/handle/star/folder", handle.Star_Folder)
		private.GET("/load/trash/folder", handle.Load_Trashes_Folders)
		private.GET("/load/trash/file", handle.Load_Trashes_Files)
		private.GET("/load/star/file", handle.Load_Stars_Files)
		private.GET("/load/shared/file", handle.LoadSharedFiles)
		private.GET("/load/star/folder", handle.Load_Stars_Folders)
		private.DELETE("/handle/delete/file", handle.Delete_Permanently_Files)
		private.DELETE("/handle/delete/folder", handle.Delete_Permanently_Folders)
		private.POST("/share/file/url", handle.ShareFileURL)
		private.POST("/share/file/gmail", handle.ShareFileGmail)
		private.GET("/find/files", handle.FindFiles)
		private.GET("/get/file/path", handle.GetPathFile)
		private.POST("/handle/rename/folder", handle.RenameThisFolder)
		private.POST("/upload/avatar", handle.UploadAvatar)
		private.POST("/change/theme", handle.ChangeTheme)
		private.DELETE("/delete/account", handle.DeleteAccount)
	}

	return engine
}
