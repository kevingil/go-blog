package trending_products_service

//google.golang.org/genproto/googleapis/analytics/data/v1beta
//https://developers.google.com/analytics/devguides/reporting/data/v1/quickstart-client-libraries#go


import (
	"context"
	"fmt"
	"log"
	"os"
	"time" 

	"google.golang.org/genproto/googleapis/analytics/data/v1beta"
)

func main() {

	ClientStart ()

	// Last 7 days
	// TODO, set by user
	reportStartRange := -7
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
	}

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
