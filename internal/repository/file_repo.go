package repository

import (
	"Stash/config"
	"fmt"
)

type Response struct {
	File_id      string
	Name         string
	ParentFolder string
	MimeType     string
	Size         int64
}

func Files_Loading(userID string) ([]Response, error) {
	var result []map[string]interface{}
	_, err := config.SupabaseAdmin.From("files").
		Select("id, name, parent_id, mime_type, size", "", false).
		Eq("owner_id", userID).
		Eq("is_trashed", "false").
		Eq("is_trashed_spread", "false").
		ExecuteTo(&result)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	Files := make([]Response, 0, len(result))
	for _, row := range result {
		Files = append(Files, Response{
			File_id:      row["id"].(string),
			Name:         row["name"].(string),
			ParentFolder: row["parent_id"].(string),
			MimeType:     row["mime_type"].(string),
			Size:         int64(row["size"].(float64)),
		})
	}
	return Files, nil
}
