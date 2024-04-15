package models

import (
	"time"

	"github.com/kevingil/blog/internal/database"
)

type File struct {
	FileId    int
	FileName  string
	FileType  string
	FileSize  int
	S3Url     string
	CreatedAt time.Time
}

func (f *File) Create() error {
	_, err := database.Db.Exec("INSERT INTO files (file_name, file_type, file_size, s3_url) VALUES (?, ?, ?, ?)", f.FileName, f.FileType, f.FileSize, f.S3Url)
	if err != nil {
		return err
	}
	return nil
}
func (f *File) Delete() error {
	_, err := database.Db.Exec("DELETE FROM files WHERE file_id = ?", f.FileId)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) Find() error {
	rows, err := database.Db.Query("SELECT file_id, file_name, file_type, file_size, s3_url, created_at FROM files WHERE file_id = ?", f.FileId)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&f.FileId, &f.FileName, &f.FileType, &f.FileSize, &f.S3Url, &f.CreatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}
