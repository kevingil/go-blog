package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// Counts page views from current day to n days ago
func CountPageViews(int n) int {
	// Load the credentials file
	credentials, err := google.ReadFile("/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Initialize the configuration
	config, err := google.JWTConfigFromJSON(credentials, analyticsdata.AnalyticsdataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// Create a client with additional options
	client := config.Client(context.Background(), option.WithScopes(analyticsdata.AnalyticsdataReadonlyScope))

	// Create the Analytics Data API service
	service, err := analyticsdata.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Error creating Analytics Data service: %v", err)
	}

	// Set the GA4 property ID
	propertyID := "ga4-property-id"

	// Get the current date and the date 30 days ago
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	// Create a request to get page views for the last 30 days
	request := &analyticsdata.RunReportRequest{
		Property: fmt.Sprintf("properties/%s", propertyID),
		DateRanges: []*analyticsdata.DateRange{
			{
				StartDate: startDate,
				EndDate:   endDate,
			},
		},
		Metrics: []*analyticsdata.Metric{
			{
				Name: "metrics/pageviews",
			},
		},
		Dimensions: []*analyticsdata.Dimension{
			{
				Name: "dimensions/pagePath",
			},
		},
	}

	// Execute the request
	response, err := service.Properties.RunReport(request).Do()
	if err != nil {
		log.Fatalf("Error executing report request: %v", err)
	}

	// Print the page views for each page path for testing
	for _, row := range response.Rows {
		pagePath := row.Dimensions[0]
		pageViews := row.Metrics[0].Values[0]
		fmt.Printf("Page Path: %s, Page Views: %s\n", pagePath, pageViews)
	}
	return nil
}
