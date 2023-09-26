package google-analytics-api


import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/analytics/data/v1beta"
	"google.golang.org/grpc"
)

func main() {
	reportStartRange := -7
	// Create context
	propertyID := os.Getenv("GA_PROPERTY_ID")
	apiKey := "GA_API_KEY"
	ctx := context.Background()
	opts := option.WithAPIKey(apiKey)

	// Data client
	client, err := data.NewAnalyticsDataClient(ctx, opts)
	if err != nil {
		log.Fatalf("Error generating Analytics Data Client: %v", err)
	}
	defer client.Close()

	// Last 7 days
	// TODO, set by user
	startDate := time.Now().AddDate(0, 0, reportStartRange).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	// Top visited query
	request := &data.RunReportRequest{
		Property: fmt.Sprintf("properties/%s", propertyID),
		DateRanges: []*data.DateRange{
			{
				StartDate: startDate,
				EndDate:   endDate,
			},
		},
		Dimensions: []*data.Dimension{
			{
				Name: "pagePath",
			},
		},
		Metrics: []*data.Metric{
			{
				Name: "activeUsers",
			},
		},
		OrderBys: []*data.OrderBy{
			{
				Metric: &data.OrderBy_Metric{
					MetricName: "activeUsers",
				},
				SortOrder: data.OrderBy_DESCENDING,
			},
		},
	};

	// Run report
	response, err := client.RunReport(ctx, request)
	if err != nil {
		log.Fatalf("Error running report, check request. %v", err)
	}

	// Print 
	for _, row := range response.Rows {
		pagePath := row.DimensionValues[0].Value
		activeUsers := row.MetricValues[0].Value
		fmt.Printf("Page Path: %s, Active Users: %v\n", pagePath, activeUsers)
	}
}

