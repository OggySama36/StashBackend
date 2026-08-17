package service

import (
	"Stash/config"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
)

func ChangePwd(userId, email, oldPwd, newPwd string) error {
	_, errCheck := config.SupabasePassword.SignInWithEmailPassword(email, oldPwd)
	if errCheck != nil {
		fmt.Println(errCheck)
		errType := errCheck.Error()
		if strings.Contains(errType, "invalid_credentials") {
			return fmt.Errorf("Your old password is incorrect!")
		}
	}

	_, errChangePassword := config.SupabasePassword.Auth.WithToken(config.App.SupabaseKey).
		AdminUpdateUser(types.AdminUpdateUserRequest{
			UserID:   uuid.MustParse(userId),
			Password: newPwd,
		})
	if errChangePassword != nil {
		fmt.Println(errChangePassword)
		return fmt.Errorf("Failed to change new password: %w", errChangePassword)
	}
	return nil
}

type GetUidAndAuthentication struct {
	UserID               string `json:"id"`
	ResetRequestVerified bool   `json:"reset_password_verified"`
}

func ResetForgotPassword(email, newPassword string) error {
	var GetUid []GetUidAndAuthentication
	_, errGet := config.SupabaseAdmin.From("profiles").
		Select("id, reset_password_verified", "", false).
		Eq("email", email).
		ExecuteTo(&GetUid)
	if errGet != nil {
		return fmt.Errorf("Failed to query user: %w", errGet)
	}
	if len(GetUid) == 0 {
		return fmt.Errorf("User not found")
	}

	if !GetUid[0].ResetRequestVerified {
		return fmt.Errorf("Your email has not been verified! Try again")
	}

	uuidUserID := uuid.MustParse(GetUid[0].UserID)
	_, errChangePassword := config.SupabasePassword.Auth.WithToken(config.App.SupabaseKey).
		AdminUpdateUser(types.AdminUpdateUserRequest{
			UserID:   uuidUserID,
			Password: newPassword,
		})
	if errChangePassword != nil {
		return fmt.Errorf("Failed to change new password: %w", errChangePassword)
	}

	_, _, errUpdateDB := config.SupabaseAdmin.From("profiles").
		Update(map[string]interface{}{
			"OTP_regain_password":           "",
			"reset_password_verified":       false,
			"OTP_regain_password_expire_at": time.Now(),
		}, "", "").
		Eq("email", email).
		Execute()
	if errUpdateDB != nil {
		return fmt.Errorf("%w", errUpdateDB)
	}
	return nil
}
