package analytics

import (
	"context"
	"strconv"
	"time"

	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
)

// Counts page views from current day to n days ago
func CountPageViews(n int) (int, error) {
	ctx := context.Background()

	// Replace 'YOUR_SERVICE_ACCOUNT_FILE.json' with the path to your service account JSON file
	analyticsService, err := analyticsdata.NewService(ctx)
	if err != nil {
		return 0, err
	}

	// Replace 'YOUR_PROPERTY_ID' with your Google Analytics property ID
	propertyID := "YOUR_PROPERTY_ID"

	// Calculate the start date and end date for the query
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().Add(-time.Duration(n) * 24 * time.Hour).Format("2006-01-02")

	// Create a request to get the page views for the specified date range
	request := &analyticsdata.RunReportRequest{
		DateRanges: []*analyticsdata.DateRange{
			{
				StartDate: startDate,
				EndDate:   endDate,
			},
		},
		Metrics: []*analyticsdata.Metric{
			{
				Expression: "ga:pageviews",
			},
		},
		Dimensions: []*analyticsdata.Dimension{
			{
				Name: "ga:pagePath",
			},
		},
	}

	// Execute the request
	response, err := analyticsService.Properties.RunReport(propertyID, request).Context(ctx).Do()
	if err != nil {
		return 0, err
	}

	// Sum up the page views from the response
	totalPageViews := 0
	for _, row := range response.Rows {
		pageViews, err := strconv.Atoi(row.ForceSendFields[0])
		if err != nil {
			print(pageViews)
			return 0, err
		}
		totalPageViews += pageViews
	}

	return totalPageViews, nil
}
