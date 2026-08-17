package service

import (
	"Stash/config"
	"fmt"

	"github.com/supabase-community/gotrue-go/types"
)

type RegisterResult struct {
	Email      string
	Name       string
	JoinedDate string
	AvatarLink string
	Theme      string
	Token      string
}

func Register(Username, Password, Email string) (*RegisterResult, error) {
	var resultQuery []map[string]interface{}
	config.SupabaseAdmin.From("profiles").
		Select("id", "", false).
		Eq("email", Email).
		ExecuteTo(&resultQuery)
	if len(resultQuery) > 0 {
		return nil, fmt.Errorf("Email has already been existed!")
	}
	saveSystem, errSave := config.SupabaseAdmin.Auth.Signup(types.SignupRequest{
		Email:    Email,
		Password: Password,
	})
	if errSave != nil {
		fmt.Println(errSave.Error())
		return nil, fmt.Errorf("Too many signup attempts, please try again later")
	}
	var newIdParsed string
	var newId []map[string]interface{}
	config.SupabaseAdmin.From("profiles").
		Insert(map[string]interface{}{
			"id":       saveSystem.User.ID,
			"username": Username,
			"email":    Email,
		}, false, "", "representation", "").
		ExecuteTo(&newId)

	if len(newId) > 0 {
		newIdParsed = newId[0]["id"].(string)
	}
	var getUserData []map[string]interface{}
	config.SupabaseAdmin.From("profiles").
		Select("create_at, avatar, theme", "", false).
		Eq("id", newIdParsed).
		Eq("email", Email).
		ExecuteTo(&getUserData)
	JoinedDate, _ := getUserData[0]["create_at"].(string)
	AvatarLink, _ := getUserData[0]["avatar"].(string)
	Theme, _ := getUserData[0]["theme"].(string)
	return &RegisterResult{
		Email:      Email,
		Name:       Username,
		JoinedDate: JoinedDate,
		AvatarLink: AvatarLink,
		Theme:      Theme,
		Token:      saveSystem.AccessToken,
	}, nil
}
