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
var TemplateFS embed.FS

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
	//users are authenticated via JWT cookies
	controllers.Sessions = make(map[string]*models.User)
	controllers.Tmpl = parseTemplates()
	controllers.Fs = TemplateFS

	// Start HTTP server
	server.Boot()

}

func parseTemplates() *template.Template {
	tmpl := template.New("").Funcs(helpers.Functions)

	dirs := []string{
		"internal/templates/pages",
		"internal/templates/partials",
	}

	for _, dir := range dirs {
		files, err := TemplateFS.ReadDir(dir)
		if err != nil {
			log.Fatalf("Error reading directory %s: %v", dir, err)
		}
		for _, file := range files {
			if file.IsDir() {
				subFiles, err := TemplateFS.ReadDir(filepath.Join(dir, file.Name()))
				if err != nil {
					log.Fatalf("Error reading directory %s: %v", filepath.Join(dir, file.Name()), err)
				}
				for _, subFile := range subFiles {
					if !subFile.IsDir() && strings.HasSuffix(subFile.Name(), ".gohtml") {
						log.Printf("Parsing file: %s", filepath.Join(dir, file.Name(), subFile.Name()))
						_, err = tmpl.ParseFS(TemplateFS, filepath.Join(dir, file.Name(), subFile.Name()))
						if err != nil {
							log.Fatalf("Error parsing file %s: %v", filepath.Join(dir, file.Name(), subFile.Name()), err)
						}
					}
				}
			} else if strings.HasSuffix(file.Name(), ".gohtml") {
				log.Printf("Parsing file: %s", filepath.Join(dir, file.Name()))
				_, err = tmpl.ParseFS(TemplateFS, filepath.Join(dir, file.Name()))
				if err != nil {
					log.Fatalf("Error parsing file %s: %v", filepath.Join(dir, file.Name()), err)
				}
			}
		}
	}

	return tmpl
}
