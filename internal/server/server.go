package server

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/kevingil/blog/internal/controllers"
	"github.com/kevingil/blog/internal/helpers"
)

func Boot() {
	// Create a new engine by passing the template folder
	// and template extension using <engine>.New(dir, ext string)
	engine := html.New("./internal/templates", ".gohtml")
	// Add your helper functions to the template's global function map.
	engine.AddFunc("until", helpers.Until)
	engine.AddFunc("date", helpers.Date)
	engine.AddFunc("shortDate", helpers.ShortDate)
	engine.AddFunc("v", helpers.V)
	engine.AddFunc("mdToHTML", helpers.MdToHTML)
	engine.AddFunc("truncate", helpers.Truncate)
	engine.AddFunc("draft", helpers.Draft)
	engine.AddFunc("ValidateEmail", helpers.ValidateEmail)
	engine.AddFunc("sub", func(a, b int) int {
		return a - b
	})
	engine.AddFunc("add", func(a, b int) int {
		return a + b
	})

	// Reload the templates on each render, good for development
	engine.Reload(true) // Optional. Default: false

	// Debug will print each template that is parsed, good for debugging
	engine.Debug(true) // Optional. Default: false

	// After you created your engine, you can pass it to Fiber's Views Engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(LayoutMiddleware)

	// User login, logout, register
	app.Get("/login", controllers.Login)
	app.Get("/logout", controllers.Logout)
	app.Get("/register", controllers.Register)

	// View posts, preview drafts
	app.Get("/blog", controllers.BlogPage)

	// Services
	app.Get("/blog/partial/recent", controllers.RecentPostsPartial)

	// View posts, preview drafts
	app.Get("/blog/:slug", controllers.BlogPostPage)

	// User Dashboard
	app.Get("/dashboard", controllers.Dashboard)

	// Edit articles, delete, or create new
	app.Get("/dashboard/articles", controllers.DashboardArticles)

	// View posts, preview drafts
	app.Get("/dashboard/articles/edit", controllers.EditArticle)

	// User Profile
	// Edit about me, skills, and projects
	app.Get("/dashboard/profile", controllers.DashboardProfile)

	// Resume Edit
	app.Get("/dashboard/resume", controllers.DashboardResume)

	// Files page
	app.Get("/dashboard/files", controllers.DashboardFilesPage)
	// Files = content with pagination
	app.Get("/dashboard/files/content", controllers.FilesContent)

	// Pages
	app.Get("/about", controllers.About)
	app.Get("/contact", controllers.Contact)

	// Combine file server and index handler
	app.Use("/", func(c *fiber.Ctx) error {
		// Serve index for root path
		if c.Path() == "/" {
			return controllers.Index(c)
		}

		// Check if the requested file exists
		path := "web" + c.Path()
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// If the file doesn't exist, serve the index page
			return controllers.Index(c)
		}

		// If the file exists, serve it
		return c.SendFile(path)
	})

	log.Printf("Your app is running on port %s.", os.Getenv("PORT"))
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
