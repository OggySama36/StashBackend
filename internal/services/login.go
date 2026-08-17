package service

import (
	"Stash/config"
	"fmt"
	"strings"
)

type LoginResult struct {
	Email      string
	Name       string
	JoinedDate string
	AvatarLink string
	Theme      string
	Token      string
}

func Login(Email, Password string) (error, *LoginResult) {
	response, errVrf := config.SupabaseAdmin.Auth.SignInWithEmailPassword(Email, Password)
	if errVrf != nil {
		errMsg := errVrf.Error()
		switch {
		case strings.Contains(errMsg, "invalid_credentials"):
			return fmt.Errorf("Your email or password is incorrect"), nil
		case strings.Contains(errMsg, "email_not_confirmed"):
			return fmt.Errorf("Your email is not confirmed"), nil
		case strings.Contains(errMsg, "user_not_found"):
			return fmt.Errorf("This account does not exist"), nil
		default:
			return fmt.Errorf("Log in false: %v", errMsg), nil
		}
	}

	userIdParse := response.User.ID.String()
	var Result []map[string]interface{}
	_, errFind := config.SupabaseAdmin.From("profiles").
		Select("username, create_at, avatar, theme", "", false).
		Eq("email", Email).
		Eq("id", userIdParse).
		ExecuteTo(&Result)

	if errFind != nil {
		fmt.Println(errFind)
		return fmt.Errorf("Failed to get name: %w", errFind), nil
	}
	if len(Result) == 0 {
		return fmt.Errorf("No username found!"), nil
	}

	NameUser, _ := Result[0]["username"].(string)
	JoinedDate, _ := Result[0]["create_at"].(string)
	AvatarLink, _ := Result[0]["avatar"].(string)
	Theme, _ := Result[0]["theme"].(string)
	return nil, &LoginResult{
		Email:      Email,
		Name:       NameUser,
		JoinedDate: JoinedDate,
		AvatarLink: AvatarLink,
		Theme:      Theme,
		Token:      response.AccessToken,
	}
}
