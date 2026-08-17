package repository

import (
	"Stash/config"
	"fmt"
)

type result_response struct {
	Folder_id    string
	Name_Folder  string
	ParentFolder string
}

func Folders_Loading(userID string) ([]result_response, error) {
	var result []map[string]interface{}
	_, err := config.SupabaseAdmin.From("folders").
		Select("id, name_folder, parent_folder", "", false).
		Eq("owner_id", userID).
		Eq("is_trashed", "false").
		Eq("is_trashed_spread", "false").
		Order("created_at", nil).
		ExecuteTo(&result)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	folders := make([]result_response, 0, len(result))
	for _, row := range result {
		folders = append(folders, result_response{
			Folder_id:    row["id"].(string),
			Name_Folder:  row["name_folder"].(string),
			ParentFolder: row["parent_folder"].(string),
		})
	}
	return folders, nil
}
