package main

import (
	"embed"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/kevingil/blog/internal/controllers"
	"github.com/kevingil/blog/internal/database"
	"github.com/kevingil/blog/internal/helpers"
	"github.com/kevingil/blog/internal/models"
	"github.com/kevingil/blog/internal/server"
)

// Tempate files for embedded file system
//
//go:embed internal/templates/pages/*.gohtml internal/templates/pages/*/*.gohtml internal/templates/partials/*.gohtml
var Fs embed.FS

// Entrypoint

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Every app instance
	// initializes a database connection
	err = database.Init()
	if err != nil {
		log.Fatal("Database initialization error:", err)
	}

	//Current session is stored in memory
	//users are authenticated via JWT tokens
	controllers.Fs = Fs
	controllers.Sessions = make(map[string]*models.User)
	controllers.Tmpl = parseTemplates(Fs)

	// Start HTTP server
	server.Boot()

}

func parseTemplates(fs embed.FS) *template.Template {
	tmpl := template.New("").Funcs(helpers.Functions)
	var parseErr error

	var walkDir func(path string) error
	walkDir = func(path string) error {
		entries, err := fs.ReadDir(path)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			entryPath := filepath.Join(path, entry.Name())
			if entry.IsDir() {
				if err := walkDir(entryPath); err != nil {
					return err
				}
			} else if filepath.Ext(entryPath) == ".gohtml" {
				// Read the content of the template file
				fileContent, err := fs.ReadFile(entryPath)
				if err != nil {
					parseErr = err
					return err
				}

				// Parse the template using its name and the read content
				templateName := strings.TrimPrefix(entryPath, "internal/templates/")
				if _, err := tmpl.New(templateName).Parse(string(fileContent)); err != nil {
					parseErr = err
					return err
				}
			}
		}
		return nil
	}

	// Walk from the root template directory
	if err := walkDir("internal/templates"); err != nil {
		log.Fatal(err)
	}
	if parseErr != nil {
		log.Fatal(parseErr)
	}

	return tmpl
}
