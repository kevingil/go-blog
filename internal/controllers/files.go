package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/kevingil/blog/pkg/storage"
)

type FilesData struct {
	Config       storage.Config
	Folders      []storage.Folder
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

	var files []storage.File
	var folders []storage.Folder

	var fileSession = storage.Session{
		UrlPrefix:       os.Getenv("CDN_URL_PREFIX"),
		BucketName:      os.Getenv("CDN_BUCKET_NAME"),
		AccountId:       os.Getenv("CDN_ACCOUNT_ID"),
		AccessKeyId:     os.Getenv("CDN_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("CDN_ACCESS_KEY_SECRET"),
		Endpoint:        os.Getenv("CDN_API_ENDPOINT"),
		Region:          "us-west-2",
	}

	fileSession, err := fileSession.Connect()
	if err != nil {
		log.Print(err)
	} else {

		files, folders, err = fileSession.List("blog", "")
		if err != nil {
			log.Print(err)
		}

	}

	filesData := FilesData{
		Files:   files,
		Folders: folders,
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
