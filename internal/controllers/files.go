package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kevingil/blog/pkg/storage"
)

var FileSession storage.Session

func AdminFilesPage(c *fiber.Ctx) error {
	data := map[string]interface{}{}
	if c.Get("HX-Request") == "true" {
		return c.Render("adminFilesPage", data, "")
	} else {
		return c.Render("adminFilesPage", data)
	}
}

func FilesContent(c *fiber.Ctx) error {
	var files []storage.File
	var folders []storage.Folder

	fileSession, err := FileSession.Connect()
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
		return c.Render("adminFilesContent", data, "")
	} else {
		return c.Render("adminFilesContent", data)
	}
}
