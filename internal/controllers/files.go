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

func HandleFileUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File upload failed"})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer src.Close()

	// Connect to the storage session
	fileSession, err := FileSession.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to connect to storage"})
	}

	// Upload the file
	err = fileSession.Upload("blog", file.Filename, src)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload file"})
	}

	return c.JSON(fiber.Map{"message": "File uploaded successfully"})
}

func UpdateDirectory(c *fiber.Ctx) error {
	currentDir := c.FormValue("currentDir")
	newDir := c.FormValue("newDir")

	// Connect to the storage session
	fileSession, err := FileSession.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to connect to storage"})
	}

	// Update the directory
	err = fileSession.UpdateFolder("blog", currentDir, newDir)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update directory"})
	}

	return c.JSON(fiber.Map{"message": "Directory updated successfully"})
}

func CreateNewDirectory(c *fiber.Ctx) error {
	newDir := c.FormValue("newDir")

	// Connect to the storage session
	fileSession, err := FileSession.Connect()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to connect to storage"})
	}

	// Create the new directory
	err = fileSession.CreateFolder("blog", newDir)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create new directory"})
	}

	return c.JSON(fiber.Map{"message": "New directory created successfully"})
}
