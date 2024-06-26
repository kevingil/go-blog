package controllers

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/kevingil/blog/pkg/storage"
)

func DashboardFilesPage(c *fiber.Ctx) error {
	data := map[string]interface{}{}
	if c.Get("HX-Request") == "true" {
		return c.Render("dashboardFilesPage", data, "")
	} else {
		return c.Render("dashboardFilesPage", data)
	}
}

func FilesContent(c *fiber.Ctx) error {
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

	data := map[string]interface{}{
		"Files":   files,
		"Folders": folders,
		"Error":   err,
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("dashboardFilesContent", data, "")
	} else {
		return c.Render("dashboardFilesContent", data)
	}
}
