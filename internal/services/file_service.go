package service

import (
	"Stash/config"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	storage_go "github.com/supabase-community/storage-go"
)

type Response struct {
	File_id  string
	Message  string
	MimeType string
	URL      string
	Size     int64
}

func File(file *multipart.FileHeader, parent_folder, userID string) (*Response, error) {
	mimeType := file.Header.Get("Content-Type")
	var newID []map[string]interface{}
	_, err := config.SupabaseAdmin.From("files").
		Insert(map[string]interface{}{
			"name":      file.Filename,
			"parent_id": parent_folder,
			"owner_id":  userID,
			"mime_type": mimeType,
			"size":      file.Size,
		}, false, "", "representation", "").
		ExecuteTo(&newID)
	if err != nil {
		fmt.Println("err: %w", err)
		return nil, err
	}
	newFilesID := newID[0]["id"].(string)
	Ext_file := filepath.Ext(file.Filename)
	StoragePath := fmt.Sprintf("%s/%s%s", userID, newFilesID, Ext_file)
	openFile, errOpen := file.Open()
	if errOpen != nil {
		fmt.Println("errOpen: %w", errOpen)
		return nil, errOpen
	}
	fileBytes, errRead := io.ReadAll(openFile)
	if errRead != nil {
		fmt.Println("errRead: %w", errRead)
		return nil, errRead
	}
	result, err := config.SupabaseAdmin.Storage.UploadFile("Stash", StoragePath, bytes.NewReader(fileBytes), storage_go.FileOptions{
		ContentType: &mimeType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload to storage: %w", err)
	}
	cleanPath := strings.TrimPrefix(result.Key, "Stash/")
	_, _, errAddURL := config.SupabaseAdmin.From("files").
		Update(map[string]interface{}{
			"storage_path": cleanPath,
		}, "", "").
		Eq("id", newFilesID).
		Execute()
	if errAddURL != nil {
		fmt.Print("errAddUrl: %w", errAddURL)
		return nil, errAddURL
	}
	return &Response{
		File_id:  newFilesID,
		Message:  "Up file ok",
		MimeType: mimeType,
		Size:     file.Size,
	}, nil
}

type Response_handle_file struct {
	Message string
}

func RemoveFile(userID, file_id string) (*Response_handle_file, error) {
	_, _, errRemove := config.SupabaseAdmin.From("files").
		Update(map[string]interface{}{
			"is_trashed": true,
			"is_starred": false,
		}, "", "").
		Eq("id", file_id).
		Execute()
	if errRemove != nil {
		return nil, fmt.Errorf("%w", errRemove)
	}
	return &Response_handle_file{
		Message: "remove ok",
	}, nil
}

func RestoreFile(userID, file_id string) (*Response_handle_file, error) {
	_, _, errRestore := config.SupabaseAdmin.From("files").
		Update(map[string]interface{}{
			"is_trashed": false,
		}, "", "").
		Eq("id", file_id).
		Eq("owner_id", userID).
		Execute()
	if errRestore != nil {
		return nil, fmt.Errorf("%w", errRestore)
	}

	return &Response_handle_file{
		Message: "restore ok",
	}, nil
}

func StarFile(userID, file_id string) (*Response_handle_file, error) {
	_, _, errRemove := config.SupabaseAdmin.From("files").
		Update(map[string]interface{}{
			"is_starred": true,
		}, "", "").
		Eq("id", file_id).
		Execute()
	if errRemove != nil {
		return nil, fmt.Errorf("%w", errRemove)
	}
	return &Response_handle_file{
		Message: "Star ok",
	}, nil
}

func UnstarFile(userID, file_id string) (*Response_handle_file, error) {
	_, _, errRemove := config.SupabaseAdmin.From("files").
		Update(map[string]interface{}{
			"is_starred": false,
		}, "", "").
		Eq("id", file_id).
		Execute()
	if errRemove != nil {
		return nil, fmt.Errorf("%w", errRemove)
	}
	return &Response_handle_file{
		Message: "Star ok",
	}, nil
}

func Delete_Files(userID, file_id string) (*Response_handle_file, error) {
	var result []map[string]interface{}
	config.SupabaseAdmin.From("files").
		Select("storage_path", "", false).
		Eq("id", file_id).
		Eq("owner_id", userID).
		ExecuteTo(&result)
	if result[0]["storage_path"] == nil {
		return nil, fmt.Errorf("Storage path not found")
	}
	_, errDeleteBucket := config.SupabaseAdmin.Storage.RemoveFile("Stash", []string{result[0]["storage_path"].(string)})
	if errDeleteBucket != nil {
		fmt.Println(errDeleteBucket)
		return nil, fmt.Errorf("%w", errDeleteBucket)
	}
	_, _, errDeleteDB := config.SupabaseAdmin.From("files").
		Delete("", "").
		Eq("id", file_id).
		Eq("owner_id", userID).
		Execute()

	if errDeleteDB != nil {
		return nil, fmt.Errorf("%w", errDeleteDB)
	}
	return &Response_handle_file{
		Message: "Delete complete",
	}, nil
}

func DeleteAllFiles(userID string) error {
	_, _, errDelete := config.SupabaseAdmin.From("files").
		Delete("", "").
		Eq("owner_id", userID).
		Execute()
	if errDelete != nil {
		fmt.Println(errDelete)
		return fmt.Errorf("Failed to delete all files: %w", errDelete)
	}
	return nil
}
