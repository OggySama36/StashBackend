package config

import "os"

type Config struct {
	AllowedOrigins      string
	SupabaseURL         string
	SupabaseKey         string
	Port                string
	ApplicationPwd      string
	HostEmail           string
	SmtpHost            string
	CloudinaryName      string
	CloudinaryAPIKey    string
	CloudinarySecretKey string
}

var App = &Config{}

func Load() {
	App.AllowedOrigins = os.Getenv("ALLOWED_ORIGINS")
	App.SupabaseURL = os.Getenv("SUPABASE_URL")
	App.SupabaseKey = os.Getenv("SUPABASE_SECRET_KEY")
	App.ApplicationPwd = os.Getenv("APPLICATION_PASSWORD")
	App.HostEmail = os.Getenv("HOST_EMAIL")
	App.SmtpHost = os.Getenv("SMTP_HOST")
	App.Port = os.Getenv("PORT")
	App.CloudinaryName = os.Getenv("CLOUDINARY_NAME")
	App.CloudinaryAPIKey = os.Getenv("CLOUDINARY_API_KEY")
	App.CloudinarySecretKey = os.Getenv("CLOUDINARY_SECRET_KEY")
}
