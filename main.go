package main

import (
	"embed"
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
//go:embed internal/templates/*.gohtml internal/templates/pages/*.gohtml internal/templates/partials/*.gohtml
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
	parsedTemplates, err := tmpl.ParseFS(templateFS,
		"internal/templates/*.gohtml",
		"internal/templates/pages/*.gohtml",
		"internal/templates/partials/*.gohtml")
	if err != nil {
		log.Fatalf("Failed to parse embedded templates: %v", err)
	}

	return parsedTemplates
}
