package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/kevingil/blog/pkg/storage"
)

type FilesData struct {
	Config       storage.Config
	Folders      []string
	Files        []storage.File
	Error        error
	TotalItems   int
	ItemsPerPage int
	TotalPages   int
	CurrentPage  int
}

func FilesPage(w http.ResponseWriter, r *http.Request) {

	req := Request{
		W:      w,
		R:      r,
		Layout: "dashboard-layout",
		Tmpl:   "dashboard-files",
		Data:   nil,
	}

	render(req)
}

func FilesContent(w http.ResponseWriter, r *http.Request) {
	s3config := storage.Config{
		AccessKey:    os.Getenv("CDN_ACCESS_KEY_ID"),
		SecretKey:    os.Getenv("CDN_SECRET_ACCESS_KEY"),
		SessionToken: os.Getenv("CDN_SESSION_TOKEN"),
		Endpoint:     os.Getenv("CDN_URL"),
		Region:       "us-west-2",
	}

	fileSession, _ := storage.NewSession(s3config)
	files, list, err := fileSession.List("blog", "")
	log.Print(err)
	filesData := FilesData{
		Files:   files,
		Folders: list,
		Error:   err,
	}

	req := Request{
		W:      w,
		R:      r,
		Layout: "dashboard-layout",
		Tmpl:   "dashboard-files-content",
		Data:   filesData,
	}

	render(req)
}
