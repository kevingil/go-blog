package controllers

import (
	"net/http"
	"os"

	"github.com/kevingil/blog/pkg/storage"
)

type FilesData struct {
	Config       storage.Config
	Folders      []string
	Files        []storage.File
	TotalItems   int
	ItemsPerPage int
	TotalPages   int
	CurrentPage  int
}

func Files(w http.ResponseWriter, r *http.Request) {
	s3config := storage.Config{
		AccessKey:    os.Getenv("CDN_ACCESS_KEY_ID"),
		SecretKey:    os.Getenv("CDN_SECRET_ACCESS_KEY"),
		SessionToken: os.Getenv("CDN_SESSION_TOKEN"),
		Endpoint:     os.Getenv("CDN_URL"),
		Region:       "us-west-2",
	}

	req := Request{
		W:      w,
		R:      r,
		Layout: "dashboard-layout",
		Tmpl:   "dashboard-files",
		Data:   nil,
	}

	fileSession, err := storage.NewSession(s3config)
	if err != nil {
		files, list, _ := fileSession.List("", "")
		filesData := FilesData{
			Files:   files,
			Folders: list,
		}
		req.Data = filesData
	}

	render(req)
}
