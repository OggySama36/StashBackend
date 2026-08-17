package service

import (
	"Stash/config"
	"fmt"
)

type Response_Found_Files struct {
	File_id      string
	Name         string
	ParentFolder string
	Size         int64
}

type Response_Found_Folders struct {
	Folder_id    string
	Name_Folder  string
	ParentFolder string
}

func FindingHome(userID, value, TypeFind string) ([]Response_Found_Files, []Response_Found_Folders, error) {
	var resultFiles []map[string]interface{}
	_, errFindFiles := config.SupabaseAdmin.From("files").
		Select("id, name, parent_id, size", "", false).
		Ilike("name", "%"+value+"%").
		Eq("owner_id", userID).
		Eq("is_trashed", "false").
		Eq("is_trashed_spread", "false").
		Limit(20, "").
		ExecuteTo(&resultFiles)
	if errFindFiles != nil {
		return nil, nil, fmt.Errorf("Failed to find files: %w", errFindFiles)
	}
	ListFoundFiles := make([]Response_Found_Files, 0, len(resultFiles))
	for _, row := range resultFiles {
		ListFoundFiles = append(ListFoundFiles, Response_Found_Files{
			File_id:      row["id"].(string),
			Name:         row["name"].(string),
			ParentFolder: row["parent_id"].(string),
			Size:         int64(row["size"].(float64)),
		})
	}

	var resultFolders []map[string]interface{}
	_, errFindFolders := config.SupabaseAdmin.From("folders").
		Select("id, name_folder, parent_folder", "", false).
		Ilike("name_folder", "%"+value+"%").
		Eq("owner_id", userID).
		Eq("is_trashed", "false").
		Eq("is_trashed_spread", "false").
		Limit(20, "").
		ExecuteTo(&resultFolders)
	if errFindFolders != nil {
		return nil, nil, fmt.Errorf("Failed to find folders: %w", errFindFolders)
	}
	ListFoundFolders := make([]Response_Found_Folders, 0, len(resultFolders))
	for _, row := range resultFolders {
		ListFoundFolders = append(ListFoundFolders, Response_Found_Folders{
			Folder_id:    row["id"].(string),
			Name_Folder:  row["name_folder"].(string),
			ParentFolder: row["parent_folder"].(string),
		})
	}

	return ListFoundFiles, ListFoundFolders, nil
}

func FindingStar(userID, value, TypeFind string) ([]Response_Found_Files, []Response_Found_Folders, error) {
	var resultFiles []map[string]interface{}
	_, errFindFiles := config.SupabaseAdmin.From("files").
		Select("id, name, parent_id, size", "", false).
		Ilike("name", "%"+value+"%").
		Eq("owner_id", userID).
		Eq("is_starred", "true").
		Eq("is_trashed", "false").
		Eq("is_trashed_spread", "false").
		Limit(20, "").
		ExecuteTo(&resultFiles)
	if errFindFiles != nil {
		return nil, nil, fmt.Errorf("Failed to find files: %w", errFindFiles)
	}
	ListFoundFiles := make([]Response_Found_Files, 0, len(resultFiles))
	for _, row := range resultFiles {
		ListFoundFiles = append(ListFoundFiles, Response_Found_Files{
			File_id:      row["id"].(string),
			Name:         row["name"].(string),
			ParentFolder: row["parent_id"].(string),
			Size:         int64(row["size"].(float64)),
		})
	}

	var resultFolders []map[string]interface{}
	_, errFindFolders := config.SupabaseAdmin.From("folders").
		Select("id, name_folder, parent_folder", "", false).
		Ilike("name_folder", "%"+value+"%").
		Eq("owner_id", userID).
		Eq("is_starred", "true").
		Eq("is_trashed", "false").
		Eq("is_trashed_spread", "false").
		Limit(20, "").
		ExecuteTo(&resultFolders)
	if errFindFolders != nil {
		return nil, nil, fmt.Errorf("Failed to find folders: %w", errFindFolders)
	}
	ListFoundFolders := make([]Response_Found_Folders, 0, len(resultFolders))
	for _, row := range resultFolders {
		ListFoundFolders = append(ListFoundFolders, Response_Found_Folders{
			Folder_id:    row["id"].(string),
			Name_Folder:  row["name_folder"].(string),
			ParentFolder: row["parent_folder"].(string),
		})
	}

	return ListFoundFiles, ListFoundFolders, nil
}

func FindingTrash(userID, value, TypeFind string) ([]Response_Found_Files, []Response_Found_Folders, error) {
	var resultFiles []map[string]interface{}
	_, errFindFiles := config.SupabaseAdmin.From("files").
		Select("id, name, parent_id, size", "", false).
		Ilike("name", "%"+value+"%").
		Eq("owner_id", userID).
		Or("is_trashed.eq.true,is_trashed_spread.eq.true", "").
		Limit(20, "").
		ExecuteTo(&resultFiles)
	if errFindFiles != nil {
		return nil, nil, fmt.Errorf("Failed to find files: %w", errFindFiles)
	}
	ListFoundFiles := make([]Response_Found_Files, 0, len(resultFiles))
	for _, row := range resultFiles {
		ListFoundFiles = append(ListFoundFiles, Response_Found_Files{
			File_id:      row["id"].(string),
			Name:         row["name"].(string),
			ParentFolder: row["parent_id"].(string),
			Size:         int64(row["size"].(float64)),
		})
	}

	var resultFolders []map[string]interface{}
	_, errFindFolders := config.SupabaseAdmin.From("folders").
		Select("id, name_folder, parent_folder", "", false).
		Ilike("name_folder", "%"+value+"%").
		Eq("owner_id", userID).
		Or("is_trashed.eq.true,is_trashed_spread.eq.true", "").
		Limit(20, "").
		ExecuteTo(&resultFolders)
	if errFindFolders != nil {
		return nil, nil, fmt.Errorf("Failed to find folders: %w", errFindFolders)
	}
	ListFoundFolders := make([]Response_Found_Folders, 0, len(resultFolders))
	for _, row := range resultFolders {
		ListFoundFolders = append(ListFoundFolders, Response_Found_Folders{
			Folder_id:    row["id"].(string),
			Name_Folder:  row["name_folder"].(string),
			ParentFolder: row["parent_folder"].(string),
		})
	}

	return ListFoundFiles, ListFoundFolders, nil
}
