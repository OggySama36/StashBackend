package service

import (
	"Stash/config"
	"encoding/json"
	"fmt"
)

type folder_result struct {
	Folder_id    string
	ParentFolder string
	Message      string
}

func Folder(owner_id, FolderName, ParentFolder string) (*folder_result, error) {
	var resultFind []map[string]interface{}
	var newFolderID string
	config.SupabaseAdmin.From("folders").
		Select("id", "", false).
		Eq("name_folder", FolderName).
		Eq("owner_id", owner_id).
		Eq("parent_folder", ParentFolder).
		ExecuteTo(&resultFind)
	if len(resultFind) > 0 {
		return nil, fmt.Errorf("This folder's name has been already existed!")
	}
	var newID []map[string]interface{}
	body, err := config.SupabaseAdmin.From("folders").
		Insert(map[string]interface{}{
			"owner_id":      owner_id,
			"name_folder":   FolderName,
			"parent_folder": ParentFolder,
		}, false, "", "representation", "").
		ExecuteTo(&newID)
	if len(newID) > 0 {
		newFolderID = newID[0]["id"].(string)
		fmt.Println("id mới:", newFolderID)
	}
	fmt.Println("body:", body)
	fmt.Println("err:", err)
	fmt.Println("newID:", newID)
	return &folder_result{
		Folder_id:    newFolderID,
		ParentFolder: ParentFolder,
		Message:      "Create new folder successfully",
	}, nil
}

type CascadeFolderInfo struct {
	Folder_id       string
	Name_Folder     string
	ParentFolder    string
	IsTrashed       bool
	IsTrashedSpread bool
}

type CascadeFileInfo struct {
	File_id         string
	Name            string
	ParentFolder    string
	MimeType        string
	Size            int64
	IsTrashed       bool
	IsTrashedSpread bool
}

type Res_Folder struct {
	Message   string
	TotalSize int64
	Folders   []CascadeFolderInfo
	Files     []CascadeFileInfo
}

type DeleteFolderResult struct {
	Message      string
	TotalSize    int64
	StoragePaths []string
}

func Rename_folder(userID, folder_id, new_name string) (*Res_Folder, error) {
	_, _, errUpdate := config.SupabaseAdmin.From("folders").
		Update(map[string]interface{}{
			"name_folder": new_name,
		}, "", "").
		Eq("id", folder_id).
		Eq("owner_id", userID).
		Execute()
	if errUpdate != nil {
		return nil, fmt.Errorf("Failed to rename: %w", errUpdate)
	}
	return &Res_Folder{
		Message: "Change name for folder successfully",
	}, nil
}

func StarFolder(userID, folder_id string) (*Res_Folder, error) {
	_, _, errRemove := config.SupabaseAdmin.From("folders").
		Update(map[string]interface{}{
			"is_starred": true,
		}, "", "").
		Eq("id", folder_id).
		Execute()
	if errRemove != nil {
		return nil, fmt.Errorf("%w", errRemove)
	}
	return &Res_Folder{
		Message: "Star ok",
	}, nil
}

func UnstarFolder(userID, folder_id string) (*Res_Folder, error) {
	_, _, errRemove := config.SupabaseAdmin.From("folders").
		Update(map[string]interface{}{
			"is_starred": false,
		}, "", "").
		Eq("id", folder_id).
		Execute()
	if errRemove != nil {
		return nil, fmt.Errorf("%w", errRemove)
	}
	return &Res_Folder{
		Message: "Star ok",
	}, nil
}

func RemoveFolder(userID, folder_id string) (*Res_Folder, error) {
	result := config.SupabaseAdmin.Rpc("remove_folder_cascade", "", map[string]interface{}{
		"target_folder_id": folder_id,
		"uid":              userID,
	})

	var cascade Res_Folder
	errParse := json.Unmarshal([]byte(result), &cascade)
	if errParse != nil {
		return nil, fmt.Errorf("failed to parse cascade result: %w", errParse)
	}

	cascade.Message = "remove ok"
	return &cascade, nil
}

func RestoreFolder(userID, folder_id string) (*Res_Folder, error) {
	result := config.SupabaseAdmin.Rpc("restore_folder_cascade", "", map[string]interface{}{
		"target_folder_id": folder_id,
		"uid":              userID,
	})

	var cascade Res_Folder
	errParse := json.Unmarshal([]byte(result), &cascade)
	if errParse != nil {
		return nil, fmt.Errorf("failed to parse cascade result: %w", errParse)
	}

	cascade.Message = "restore ok"
	return &cascade, nil
}

func Delete_Folders(userID, folder_id string) error {
	result := config.SupabaseAdmin.Rpc("delete_folder_cascade", "", map[string]interface{}{
		"target_folder_id": folder_id,
		"uid":              userID,
	})
	var cascade DeleteFolderResult
	errParse := json.Unmarshal([]byte(result), &cascade)
	if errParse != nil {
		return fmt.Errorf("failed to parse cascade result: %w", errParse)
	}
	if len(cascade.StoragePaths) > 0 {
		_, errRemoveBucket := config.SupabaseAdmin.Storage.RemoveFile("Stash", cascade.StoragePaths)
		if errRemoveBucket != nil {
			fmt.Println(errRemoveBucket)
			return fmt.Errorf("Failed to delete files cascade: %w", errRemoveBucket)
		}
	}
	return nil
}

func DeleteAllFolders(userID string) error {
	_, _, errDelete := config.SupabaseAdmin.From("folders").
		Delete("", "").
		Eq("owner_id", userID).
		Execute()
	if errDelete != nil {
		fmt.Println(errDelete)
		return fmt.Errorf("Failed to delete all folders: %w", errDelete)
	}
	return nil
}
