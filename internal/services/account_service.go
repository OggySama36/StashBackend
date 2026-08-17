package service

import (
	"Stash/config"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
)

type URLAvatar struct {
	AvatarLink string
}

func SaveUrlAvatar(userID, Avatar string) (*URLAvatar, error) {
	var AvatarParsed string
	var Result []map[string]interface{}
	_, errUpdate := config.SupabaseAdmin.From("profiles").
		Update(map[string]interface{}{
			"avatar": Avatar,
		}, "", "").
		Eq("id", userID).
		ExecuteTo(&Result)
	if errUpdate != nil {
		fmt.Println(errUpdate)
		return nil, fmt.Errorf("Failed to save database: %w", errUpdate)
	}
	if len(Result) > 0 {
		AvatarParsed = Result[0]["avatar"].(string)
	}
	return &URLAvatar{
		AvatarLink: AvatarParsed,
	}, nil
}

func AnotherTheme(userID, theme string) error {
	_, _, errUpdate := config.SupabaseAdmin.From("profiles").
		Update(map[string]interface{}{
			"theme": theme,
		}, "", "").
		Eq("id", userID).
		Execute()
	if errUpdate != nil {
		fmt.Println(errUpdate)
		return fmt.Errorf("Failed to change theme: %w", errUpdate)
	}
	return nil
}

func DeleteProfile(userID string) error {
	_, _, errDelete := config.SupabaseAdmin.From("profiles").
		Delete("", "").
		Eq("id", userID).
		Execute()
	if errDelete != nil {
		fmt.Println(errDelete)
		return fmt.Errorf("Failed to delete account: %w", errDelete)
	}
	return nil
}

type OTP_code struct {
	OTP int
}

func SaveTemporaryEmail(userID, newEmail string) (*OTP_code, error) {
	var ResultQuery []map[string]interface{}
	config.SupabaseAdmin.From("profiles").
		Select("id", "", false).
		Eq("email", newEmail).
		ExecuteTo(&ResultQuery)
	if len(ResultQuery) > 0 {
		return nil, fmt.Errorf("This email has already been used by another user!")
	}
	otpCode := rand.Intn(900000) + 100000
	otpExpireAt := time.Now().Add(5 * time.Minute)
	_, _, errSaveTemporary := config.SupabaseAdmin.From("profiles").
		Update(map[string]interface{}{
			"temporary_email":            newEmail,
			"OTP_change_email":           otpCode,
			"OTP_change_email_expire_at": otpExpireAt,
		}, "", "").
		Eq("id", userID).
		Execute()
	if errSaveTemporary != nil {
		fmt.Println(errSaveTemporary)
		return nil, fmt.Errorf("Failed to save temporarily email in database: %w", errSaveTemporary)
	}
	return &OTP_code{
		OTP: otpCode,
	}, nil
}

type NewEmail struct {
	NewEmail string
}
type OTPCheckResult struct {
	TemporaryEmail         string    `json:"temporary_email"`
	OTPChangeEmail         string    `json:"OTP_change_email"`
	OTPChangeEmailExpireAt time.Time `json:"OTP_change_email_expire_at"`
}

func CheckVerifyCodeAndDecide(Otp int, userID string) (*NewEmail, error) {
	var Result []OTPCheckResult
	_, errGetValueToCheck := config.SupabaseAdmin.From("profiles").
		Select("temporary_email, OTP_change_email, OTP_change_email_expire_at", "", false).
		Eq("id", userID).
		ExecuteTo(&Result)
	if errGetValueToCheck != nil {
		fmt.Println(errGetValueToCheck)
		return nil, fmt.Errorf("Failed to check verify code: %w", errGetValueToCheck)
	}
	if len(Result) == 0 {
		return nil, fmt.Errorf("no record found")
	}
	OTP_change_email_expire_at := Result[0].OTPChangeEmailExpireAt
	temporary_email := Result[0].TemporaryEmail
	OTP_change_email, _ := strconv.Atoi(Result[0].OTPChangeEmail)
	if time.Now().After(OTP_change_email_expire_at) {
		return nil, fmt.Errorf("Your code has already been expired! Try later")
	}
	if Otp != OTP_change_email {
		return nil, fmt.Errorf("OTP code incorrect!")
	}

	_, _, errUpdateDB := config.SupabaseAdmin.From("profiles").
		Update(map[string]interface{}{
			"email":            temporary_email,
			"temporary_email":  "",
			"OTP_change_email": "",
		}, "", "").
		Eq("id", userID).
		Execute()
	if errUpdateDB != nil {
		fmt.Println(errUpdateDB)
		return nil, fmt.Errorf("Failed to change email in database: %w", errUpdateDB)
	}
	userIDuuid, _ := uuid.Parse(userID)
	_, errChangeEmail := config.SupabaseEmail.Auth.AdminUpdateUser(types.AdminUpdateUserRequest{
		UserID: userIDuuid,
		Email:  temporary_email,
	})
	if errChangeEmail != nil {
		fmt.Println(errChangeEmail)
		return nil, fmt.Errorf("Failed to change email in system: %w", errChangeEmail)
	}
	return &NewEmail{
		NewEmail: temporary_email,
	}, nil
}

func FindForgotEmail(email string) (*OTP_code, error) {
	var Result []map[string]interface{}
	config.SupabaseAdmin.From("profiles").
		Select("id", "", false).
		Eq("email", email).
		ExecuteTo(&Result)
	if len(Result) == 0 {
		return nil, fmt.Errorf("Could not find your email, please try again with another one!")
	}
	userID := Result[0]["id"].(string)
	otpCode := rand.Intn(900000) + 100000
	otpCodeExpireAt := time.Now().Add(5 * time.Minute)
	config.SupabaseAdmin.From("profiles").
		Update(map[string]interface{}{
			"OTP_regain_password":           otpCode,
			"OTP_regain_password_expire_at": otpCodeExpireAt,
		}, "", "").
		Eq("id", userID).
		Execute()
	return &OTP_code{
		OTP: otpCode,
	}, nil
}

type OTPForgotPasswordCheckResult struct {
	OTPRegainPassword         string    `json:"OTP_regain_password"`
	OTPRegainPasswordExpireAt time.Time `json:"OTP_regain_password_expire_at"`
}

func CheckCodeForgotPassword(otpCode, email string) error {
	var Result []OTPForgotPasswordCheckResult
	config.SupabaseAdmin.From("profiles").
		Select("OTP_regain_password, OTP_regain_password_expire_at", "", false).
		Eq("email", email).
		ExecuteTo(&Result)
	OTP_regain_password := Result[0].OTPRegainPassword
	OTP_regain_password_expire_at := Result[0].OTPRegainPasswordExpireAt
	if time.Now().After(OTP_regain_password_expire_at) {
		_, _, errUpdateDB := config.SupabaseAdmin.From("profiles").
			Update(map[string]interface{}{
				"OTP_regain_password":           "",
				"OTP_regain_password_expire_at": time.Now(),
			}, "", "").
			Eq("email", email).
			Execute()
		if errUpdateDB != nil {
			fmt.Println(errUpdateDB)
			return fmt.Errorf("%w", errUpdateDB)
		}
		return fmt.Errorf("Your OTP code has already been expired! Try later")
	}
	if otpCode != OTP_regain_password {
		return fmt.Errorf("OTP code incorrect!")
	}
	_, _, errUpdateDB := config.SupabaseAdmin.From("profiles").
		Update(map[string]interface{}{
			"OTP_regain_password":           "",
			"reset_password_verified":       "TRUE",
			"OTP_regain_password_expire_at": time.Now(),
		}, "", "").
		Eq("email", email).
		Execute()
	if errUpdateDB != nil {
		fmt.Println(errUpdateDB)
		return fmt.Errorf("%w", errUpdateDB)
	}
	return nil
}
