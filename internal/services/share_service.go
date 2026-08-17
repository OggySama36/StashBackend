package service

import (
	"Stash/config"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/gomail.v2"
)

type ResponseShare struct {
	Share_id string
	File_id  string
	Name     string
	MimeType string
	SharedBy string
	SharedTo string
	Method   string
}

func Copy_url_share(file_id, sender, userID string) (*ResponseShare, error) {
	var Result []map[string]interface{}
	_, err := config.SupabaseAdmin.From("files").
		Select("name, mime_type", "", false).
		Eq("id", file_id).
		Eq("owner_id", userID).
		ExecuteTo(&Result)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if len(Result) == 0 {
		return nil, fmt.Errorf("File not found or you don't have permission")
	}
	Name := Result[0]["name"].(string)
	MimeType := Result[0]["mime_type"].(string)
	var newID []map[string]interface{}
	_, errSave := config.SupabaseAdmin.From("share_logs").
		Insert(map[string]interface{}{
			"file_id":   file_id,
			"shared_by": sender,
			"shared_to": "",
			"method":    "url",
		}, false, "", "representation", "").
		ExecuteTo(&newID)
	if errSave != nil {
		fmt.Println(errSave)
		return nil, fmt.Errorf("%w", errSave)
	}
	Share_id := newID[0]["id"].(string)
	return &ResponseShare{
		Share_id: Share_id,
		File_id:  file_id,
		Name:     Name,
		MimeType: MimeType,
		SharedBy: sender,
		SharedTo: "",
		Method:   "url",
	}, nil
}

func Send_gmail_file(file_id, fileName, note, sender, recipient, userID string) (*ResponseShare, error) {
	var Result []map[string]interface{}
	_, err := config.SupabaseAdmin.From("files").
		Select("name, mime_type", "", false).
		Eq("id", file_id).
		Eq("owner_id", userID).
		ExecuteTo(&Result)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	if len(Result) == 0 {
		return nil, fmt.Errorf("File not found or you don't have permission")
	}
	fileURL, errCall := Get_Url(userID, file_id)
	if errCall != nil {
		fmt.Println(errCall)
		return nil, fmt.Errorf("Failed to get url: %w", errCall)
	}
	fileData, errGet := http.Get(fileURL.URL)
	if errGet != nil {
		return nil, fmt.Errorf("Failed to fetch file: %w", errGet)
	}
	defer fileData.Body.Close()
	m := gomail.NewMessage()
	if fileData.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch file, status: %d", fileData.StatusCode)
	}
	m.SetHeader("From", config.App.HostEmail)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "Stash sharing file...")
	m.SetBody("text/plain", note)
	m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := io.Copy(w, fileData.Body)
		return err
	}))
	send := gomail.NewDialer(config.App.SmtpHost, 587, config.App.HostEmail, config.App.ApplicationPwd)
	errSend := send.DialAndSend(m)
	if errSend != nil {
		fmt.Println(errSend)
		return nil, fmt.Errorf("Failed to send email: %w", errSend)
	}
	var newID []map[string]interface{}
	_, errSave := config.SupabaseAdmin.From("share_logs").
		Insert(map[string]interface{}{
			"file_id":   file_id,
			"shared_by": sender,
			"shared_to": recipient,
			"method":    "gmail",
		}, false, "", "representation", "").
		ExecuteTo(&newID)
	if len(newID) == 0 {
		return nil, fmt.Errorf("Failed to get new share log")
	}
	if errSave != nil {
		fmt.Println(errSave)
		return nil, fmt.Errorf("%w", errSave)
	}
	Share_id := newID[0]["id"].(string)
	Name := Result[0]["name"].(string)
	MimeType := Result[0]["mime_type"].(string)

	return &ResponseShare{
		Share_id: Share_id,
		File_id:  file_id,
		Name:     Name,
		MimeType: MimeType,
		SharedBy: sender,
		SharedTo: recipient,
		Method:   "gmail",
	}, nil
}

func LoadShared(userID, email string) ([]ResponseShare, error) {
	var Result []map[string]interface{}
	_, errQuery := config.SupabaseAdmin.From("share_logs").
		Select("id, file_id, shared_by, shared_to, method", "", false).
		Or(fmt.Sprintf("shared_by.eq.%s,shared_to.eq.%s", email, email), "").
		ExecuteTo(&Result)
	if errQuery != nil {
		fmt.Println(errQuery)
		return nil, fmt.Errorf("Failed to get shared files: %w", errQuery)
	}
	ResponseShareList := make([]ResponseShare, 0, len(Result))
	for _, row := range Result {
		file_id := row["file_id"].(string)
		var fileInfo []map[string]interface{}
		config.SupabaseAdmin.From("files").
			Select("name, mime_type", "", false).
			Eq("id", file_id).
			ExecuteTo(&fileInfo)
		if len(fileInfo) == 0 {
			return nil, fmt.Errorf("Not found file info")
		}
		name := fileInfo[0]["name"].(string)
		MimeType := fileInfo[0]["mime_type"].(string)
		shareID, _ := row["id"].(string)
		sharedBy, _ := row["shared_by"].(string)
		sharedTo, _ := row["shared_to"].(string)
		method, _ := row["method"].(string)
		ResponseShareList = append(ResponseShareList, ResponseShare{
			Share_id: shareID,
			File_id:  file_id,
			Name:     name,
			MimeType: MimeType,
			SharedBy: sharedBy,
			SharedTo: sharedTo,
			Method:   method,
		})
	}
	return ResponseShareList, nil
}
