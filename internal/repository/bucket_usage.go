package repository

import (
	"Stash/config"
)

type Response_usage struct {
	Used int64
	Left int64
}

func Get_Usage(userID string) (*Response_usage, error) {
	var result []map[string]interface{}
	var maxSize int64 = 100 * 1024 * 1024
	var used int64 = 0
	config.SupabaseAdmin.From("files").
		Select("size", "", false).
		Eq("owner_id", userID).
		Eq("is_trashed", "false").
		ExecuteTo(&result)
	if len(result) == 0 {
		return &Response_usage{
			Used: used,
			Left: maxSize,
		}, nil
	}
	for _, row := range result {
		size := row["size"].(float64)
		used += int64(size)
	}
	var Left int64 = maxSize - used
	if Left < 0 {
		Left = 0
	}
	return &Response_usage{
		Used: used,
		Left: Left,
	}, nil
}
