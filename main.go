package main

import (
	"embed"
	"io/fs"
	"log"
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
//go:embed internal/templates/pages/*.gohtml internal/templates/partials/*.gohtml
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

	// Parse templates
	controllers.Tmpl = ParseTemplates(TemplateFS)

	// Start HTTP server
	server.Boot()

}

func ParseTemplates(templateFS embed.FS) *template.Template {
	// Create a new template and add helper functions
	tmpl := template.New("").Funcs(helpers.Functions)

	// Parse the templates from the embedded file system
	pages, err := fs.Glob(templateFS, "internal/templates/pages/*.gohtml")
	if err != nil {
		log.Fatalf("Failed to read pages templates: %v", err)
	}

	partials, err := fs.Glob(templateFS, "internal/templates/partials/*.gohtml")
	if err != nil {
		log.Fatalf("Failed to read partials templates: %v", err)
	}

	// Combine the slices
	templates := append(pages, partials...)

	for _, file := range templates {
		content, err := templateFS.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read template file: %v", err)
		}

		_, err = tmpl.New(file).Parse(string(content))
		if err != nil {
			log.Fatalf("Failed to parse template: %v", err)
		}
	}

	return tmpl
}
