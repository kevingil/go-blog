package main

import (
	"embed"
	"fmt"
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
//go:embed internal/templates/pages/*.gohtml
//go:embed internal/templates/pages/**/*.gohtml
//go:embed internal/templates/pages/**/**/*.gohtml
//go:embed internal/templates/partials/*.gohtml
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

	var walkDir func(string) error
	walkDir = func(dir string) error {
		log.Printf("Entering directory: %s", dir)
		entries, err := fs.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to read directory %s: %w", dir, err)
		}

		for _, entry := range entries {
			path := filepath.Join(dir, entry.Name())
			if entry.IsDir() {
				log.Printf("Found subdirectory: %s", path)
				if err := walkDir(path); err != nil {
					return err
				}
			} else if filepath.Ext(path) == ".gohtml" {
				log.Printf("Parsing file: %s", path)
				fileContent, err := fs.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", path, err)
				}

				templateName := strings.TrimPrefix(path, "internal/templates/")
				_, err = tmpl.New(templateName).Parse(string(fileContent))
				if err != nil {
					return fmt.Errorf("failed to parse template %s: %w", templateName, err)
				}
				log.Printf("Successfully parsed template: %s", templateName)
			}
		}
		return nil
	}

	if err := walkDir("internal/templates"); err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	return tmpl
}
