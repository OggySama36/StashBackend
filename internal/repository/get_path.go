package repository

import (
	"Stash/config"
	"encoding/json"
	"fmt"
)

type PathFolder struct {
	Folder_id   string
	Name_Folder string
	Depth       int
}

func GetFolderPath(folderID string, userID string) ([]PathFolder, error) {
	if folderID == "" || folderID == "Root" {
		return []PathFolder{}, nil
	}
	result := config.SupabaseAdmin.Rpc("get_folder_path", "", map[string]interface{}{
		"start_folder_id": folderID,
		"uid":             userID,
	})
	var path []PathFolder
	if errUnmarshal := json.Unmarshal([]byte(result), &path); errUnmarshal != nil {
		return nil, fmt.Errorf("failed to parse folder path: %w", errUnmarshal)
	}
	fmt.Println("RPC raw result:", path)

	return path, nil
}
