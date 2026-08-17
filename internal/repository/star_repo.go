package repository

import (
	"Stash/config"
	"fmt"
)

type Response_star_files struct {
	File_id      string
	Name         string
	ParentFolder string
	MimeType     string
	Size         int64
}

type Response_star_folders struct {
	Folder_id    string
	Name_Folder  string
	ParentFolder string
}

func Stars_Files_Loading(userID string) ([]Response_star_files, error) {
	var result []map[string]interface{}
	_, err := config.SupabaseAdmin.From("files").
		Select("id, name, parent_id, mime_type, size", "", false).
		Eq("owner_id", userID).
		Eq("is_starred", "true").
		ExecuteTo(&result)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	Files := make([]Response_star_files, 0, len(result))
	for _, row := range result {
		Files = append(Files, Response_star_files{
			File_id:      row["id"].(string),
			Name:         row["name"].(string),
			ParentFolder: row["parent_id"].(string),
			MimeType:     row["mime_type"].(string),
			Size:         int64(row["size"].(float64)),
		})
	}
	return Files, nil
}

func Stars_Folders_Loading(userID string) ([]Response_star_folders, error) {
	var result []map[string]interface{}
	_, err := config.SupabaseAdmin.From("folders").
		Select("id, name_folder, parent_folder", "", false).
		Eq("owner_id", userID).
		Eq("is_starred", "true").
		ExecuteTo(&result)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	Folders := make([]Response_star_folders, 0, len(result))
	for _, row := range result {
		Folders = append(Folders, Response_star_folders{
			Folder_id:    row["id"].(string),
			Name_Folder:  row["name_folder"].(string),
			ParentFolder: row["parent_folder"].(string),
		})
	}
	return Folders, nil
}
