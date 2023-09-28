package trending_products_service

import (
	"context"
	"fmt"

	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"
)

// Google Analytics Property
type property struct {
	id  string
	key string
}

// TODO
// API Ready Response
// functions currently printf
// formatting is pending
type formattedResponse struct {
	response string
}

// Instance returns an Analytics Data service
func (p property) Instance() *analyticsdata.Service {

	//Analytics Data client
	ctx := context.Background()
	client, err := analyticsdata.NewService(ctx, option.WithAPIKey(p.key))
	if err != nil {
		panic(err)
	}

	return client
}

func (p property) getPopularProductsByPurchaseAmount() *formattedResponse {

	// Run the report for top items purchased by item name.
	topItemsPurchasedReportRequest := &analyticsdata.RunReportRequest{
		Property: "properties/PROPERTY_ID",
		DateRanges: []*analyticsdata.DateRange{
			{StartDate: "2023-09-27", EndDate: "2023-09-27"},
		},
		Dimensions: []*analyticsdata.Dimension{
			{Name: "item_name"},
		},
		Metrics: []*analyticsdata.Metric{
			{Name: "total_item_revenue"},
		},
		Limit: 10,
	}

	topItemsPurchasedReportResponse, err := p.Instance().Properties.RunReport(p.id, topItemsPurchasedReportRequest).Do()
	if err != nil {
		panic(err)
	}

	fmt.Println("Top items purchased by item name:")
	for _, row := range topItemsPurchasedReportResponse.Rows {
		fmt.Printf("%s: %f\n", row.DimensionValues[0].Value, row.MetricValues[0].Value)
	}
}

func (p property) getPopularProductsByViews() *formattedResponse {

	// Run the report for top views by page title and screen class.
	topViewsReportRequest := &analyticsdata.RunReportRequest{
		Property: "properties/YOUR_PROPERTY_ID",
		DateRanges: []*analyticsdata.DateRange{
			{StartDate: "2023-09-27", EndDate: "2023-09-27"},
		},
		Dimensions: []*analyticsdata.Dimension{
			{Name: "page_title"},
			{Name: "screen_class"},
		},
		Metrics: []*analyticsdata.Metric{
			{Name: "total_views"},
		},
		Limit: 10,
	}

	topViewsReportResponse, err := p.Instance().Properties.RunReport(p.id, topViewsReportRequest).Do()
	if err != nil {
		panic(err)
	}

	fmt.Println("Top views by page title and screen class:")
	for _, row := range topViewsReportResponse.Rows {
		fmt.Printf("%s (%s): %d\n", row.DimensionValues[0].Value, row.DimensionValues[1].Value, row.MetricValues[0].Value)
	}
}
