package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"github.com/kevingil/blog/internal/controllers"
	"github.com/kevingil/blog/internal/database"
	"github.com/kevingil/blog/internal/server"
	"github.com/kevingil/blog/pkg/storage"
)

// Entrypoint
func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Initializes a database connection
	err = database.Init()
	if err != nil {
		log.Fatal("Database initialization error:", err)
	}

	//Google Analytics
	controllers.AnalyticsPropertyID = os.Getenv("GA_PROPERTYID")
	controllers.AnalyticsServiceAccountJsonPath = os.Getenv("GA_SERVICE_ACCOUNT_JSON_PATH")
	//S3/R2
	controllers.FileSession = storage.Session{
		UrlPrefix:       os.Getenv("CDN_URL_PREFIX"),
		BucketName:      os.Getenv("CDN_BUCKET_NAME"),
		AccountId:       os.Getenv("CDN_ACCOUNT_ID"),
		AccessKeyId:     os.Getenv("CDN_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("CDN_ACCESS_KEY_SECRET"),
		Endpoint:        os.Getenv("CDN_API_ENDPOINT"),
		Region:          "us-west-2",
	}

	// Start HTTP server
	controllers.Store = session.New()
	server.Serve()

}
