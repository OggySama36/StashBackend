package repository

import (
	"Stash/config"
	"fmt"
)

type Response_trash_files struct {
	File_id         string
	Name            string
	ParentFolder    string
	MimeType        string
	Size            int64
	IsTrashed       bool
	IsTrashedSpread bool
}

type Response_trash_folders struct {
	Folder_id       string
	Name_Folder     string
	ParentFolder    string
	IsTrashed       bool
	IsTrashedSpread bool
}

func Trashes_Files_Loading(userID string) ([]Response_trash_files, error) {
	var result []map[string]interface{}
	_, err := config.SupabaseAdmin.From("files").
		Select("id, name, parent_id, mime_type, size, is_trashed, is_trashed_spread", "", false).
		Eq("owner_id", userID).
		Or("is_trashed.eq.true,is_trashed_spread.eq.true", "").
		ExecuteTo(&result)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	Files := make([]Response_trash_files, 0, len(result))
	for _, row := range result {
		Files = append(Files, Response_trash_files{
			File_id:         row["id"].(string),
			Name:            row["name"].(string),
			ParentFolder:    row["parent_id"].(string),
			MimeType:        row["mime_type"].(string),
			Size:            int64(row["size"].(float64)),
			IsTrashed:       row["is_trashed"].(bool),
			IsTrashedSpread: row["is_trashed_spread"].(bool),
		})
	}
	return Files, nil
}

func Trashes_Folders_Loading(userID string) ([]Response_trash_folders, error) {
	var result []map[string]interface{}
	_, err := config.SupabaseAdmin.From("folders").
		Select("id, name_folder, parent_folder, is_trashed, is_trashed_spread", "", false).
		Eq("owner_id", userID).
		Or("is_trashed.eq.true,is_trashed_spread.eq.true", "").
		ExecuteTo(&result)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	Folders := make([]Response_trash_folders, 0, len(result))
	for _, row := range result {
		Folders = append(Folders, Response_trash_folders{
			Folder_id:       row["id"].(string),
			Name_Folder:     row["name_folder"].(string),
			ParentFolder:    row["parent_folder"].(string),
			IsTrashed:       row["is_trashed"].(bool),
			IsTrashedSpread: row["is_trashed_spread"].(bool),
		})
	}
	return Folders, nil
}
