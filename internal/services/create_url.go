package service

import (
	"Stash/config"
	"fmt"
)

type Response_url struct {
	File_id string
	URL     string
}

func Get_Url(userID, File_id string) (*Response_url, error) {
	var result []map[string]interface{}
	config.SupabaseAdmin.From("files").
		Select("storage_path, owner_id", "", false).
		Eq("id", File_id).
		ExecuteTo(&result)
	fmt.Println("storagePath từ DB:", result[0]["storage_path"].(string))
	signedURL, errSign := config.SupabaseAdmin.Storage.CreateSignedUrl("Stash", result[0]["storage_path"].(string), 259200)
	if errSign != nil {
		return nil, fmt.Errorf("failed to create signed url: %w", errSign)
	}
	return &Response_url{
		File_id: File_id,
		URL:     signedURL.SignedURL,
	}, nil
}
