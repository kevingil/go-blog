package controllers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gofiber/fiber/v2"
	data "google.golang.org/genproto/googleapis/analytics/data/v1beta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var AnalyticsPropertyID string
var AnalyticsServiceAccountKeyPath string

func getAnalyticsData(c *fiber.Ctx, req *data.RunReportRequest) (*data.RunReportResponse, error) {
	_, err := GetUser(c)
	if err != nil {
		return nil, fmt.Errorf("unauthorized request")
	}

	// Load service account credentials from JSON key file
	creds, err := oauth.NewServiceAccountFromFile(AnalyticsServiceAccountKeyPath, "https://www.googleapis.com/auth/analytics.readonly")
	if err != nil {
		log.Fatalf("Failed to load credentials: %v", err)

	}

	// Create a gRPC connection with credentials
	conn, err := grpc.Dial(
		"analyticsdata.googleapis.com:443",
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
		grpc.WithPerRPCCredentials(creds),
	)
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
	}
	defer conn.Close()

	// Create a client for the Analytics Data API
	client := data.NewBetaAnalyticsDataClient(conn)

	response, err := client.RunReport(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to run report: %v", err)
		return nil, err
	}

	return response, nil
}

func GetSiteVisits(c *fiber.Ctx) error {
	scope := c.Query("range")

	var startDate string
	var endDate string

	switch {
	case scope == "all":
		// Range for all time
		startDate = "2020-01-01"
		endDate = time.Now().Format("2006-01-02")
	case strings.HasSuffix(scope, "d"):
		// Set startDate to N days ago
		days, _ := strconv.Atoi(strings.TrimSuffix(scope, "d"))
		startDate = time.Now().AddDate(0, 0, -days).Format("2006-01-02")
		endDate = time.Now().Format("2006-01-02")
	case strings.HasSuffix(scope, "mo"):
		// Set startDate N months ago
		months, _ := strconv.Atoi(strings.TrimSuffix(scope, "mo"))
		startDate = time.Now().AddDate(0, -months, 0).Format("2006-01-02")
		endDate = time.Now().Format("2006-01-02")
	default:
		// If not right format
		return c.Status(fiber.StatusBadRequest).SendString("Invalid range")
	}

	// Analyitics report request using parsed data
	req := &data.RunReportRequest{
		Property: "properties/" + AnalyticsPropertyID,
		DateRanges: []*data.DateRange{
			{
				StartDate: startDate,
				EndDate:   endDate,
			},
		},
		Dimensions: []*data.Dimension{
			{
				Name: "date",
			},
		},
		Metrics: []*data.Metric{
			{
				Name: "totalUsers",
			},
		},
	}

	response, err := getAnalyticsData(c, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to get analytics: %s", err.Error()))
	}
	var visits int
	for _, row := range response.Rows {
		userCount, _ := strconv.Atoi(row.MetricValues[0].GetValue())
		visits += userCount
	}

	return c.SendString(fmt.Sprintf("%d", visits))
}

func GetSiteVisitsChart(c *fiber.Ctx) error {
	// Calculate date range for last 6 months
	endDate := time.Now()
	startDate := endDate.AddDate(0, -6, 0)

	// Create analytics request
	req := &data.RunReportRequest{
		Property: "properties/" + AnalyticsPropertyID,
		DateRanges: []*data.DateRange{
			{
				StartDate: startDate.Format("2006-01-02"),
				EndDate:   endDate.Format("2006-01-02"),
			},
		},
		Dimensions: []*data.Dimension{
			{
				Name: "date",
			},
		},
		Metrics: []*data.Metric{
			{
				Name: "totalUsers",
			},
		},
	}

	response, err := getAnalyticsData(c, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to get analytics: %s", err.Error()))
	}

	// Process the data for the chart
	type dataPoint struct {
		Date   time.Time
		Visits int
	}
	var dataPoints []dataPoint

	for _, row := range response.Rows {
		date, err := time.Parse("20060102", row.DimensionValues[0].GetValue())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to parse date: %s", err.Error()))
		}
		visits, err := strconv.Atoi(row.MetricValues[0].GetValue())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to parse visits: %s", err.Error()))
		}
		dataPoints = append(dataPoints, dataPoint{Date: date, Visits: visits})
	}

	// Sort the data points by date
	sort.Slice(dataPoints, func(i, j int) bool {
		return dataPoints[i].Date.Before(dataPoints[j].Date)
	})

	// Compress data into 12 bi-weekly points
	compressedData := make([]dataPoint, 12)
	interval := time.Hour * 24 * 14 // 2 weeks
	currentDate := startDate
	for i := 0; i < 12; i++ {
		endInterval := currentDate.Add(interval)
		totalVisits := 0
		count := 0
		for _, dp := range dataPoints {
			if dp.Date.After(currentDate) && dp.Date.Before(endInterval) {
				totalVisits += dp.Visits
				count++
			}
		}
		averageVisits := 0
		if count > 0 {
			averageVisits = totalVisits / count
		}
		compressedData[i] = dataPoint{Date: currentDate, Visits: averageVisits}
		currentDate = endInterval
	}

	// Create slices for dates and visits
	var dates []string
	var visits []string
	for _, dp := range compressedData {
		dates = append(dates, fmt.Sprintf(`"%s",`, dp.Date.Format("2006-01-02")))
		visits = append(visits, strconv.Itoa(dp.Visits)+",")
	}

	// Create a template with the chart data
	tmpl := `
    <div id="dash-chart-container" style="height: 200px; width: 100%;">
        <canvas id="visitsChart"></canvas>
    </div>
    <script>
        new Chart(document.getElementById('visitsChart'), {
            type: 'line',
            data: {
                labels: {{.Dates}},
                datasets: [{
                    label: 'Site Visits',
                    data: {{.Visits}},
                    borderColor: 'rgb(162, 125, 255)',
                    tension: 0.1
                }]
            },
            options: {
				responsive: true,
				maintainAspectRatio: false,
				plugins: {
					legend: {
						display: false
					}
				},
                scales: {
                    y: {
                        beginAtZero: true,
						display: true
                    },
					x: {
						display: true
					}
                }
            }
        });
    </script>
    `

	// Create a template and execute it with the data
	t, err := template.New("chart").Parse(tmpl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse template")
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, struct {
		Dates  []string
		Visits []string
	}{
		Dates:  dates,
		Visits: visits,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to execute template")
	}

	// Set the content type to HTML and send the response
	c.Set("Content-Type", "text/html")
	return c.SendString(buf.String())
}

/*
	// REQUEST EXAMPLE
	// Make a request to the GA4 Data API
	request := &data.RunReportRequest{
		Property: "properties/" + AnalyticsPropertyID,
		Dimensions: []*data.Dimension{
			{Name: "eventName"},
		},
		Metrics: []*data.Metric{
			{Name: "eventCount"},
			{Name: "totalUsers"},
		},
		DateRanges: []*data.DateRange{
			{
				StartDate: time.Now().AddDate(0, 0, -28).Format("2006-01-02"),
				EndDate:   time.Now().Format("2006-01-02"),
			},
		},
	}

	// with filter
	dimensionFilter := &data.Filter{
			FieldName: "eventName",
			OneFilter: &data.Filter_StringFilter_{
				StringFilter: &data.Filter_StringFilter{
					MatchType:     data.Filter_StringFilter_EXACT,
					Value:         "banner click",
					CaseSensitive: false,
				},
			},
		}


    // Make a request to the GA4 Data API
    request := &data.RunReportRequest{
        Property: "properties/" + propertyID,
        Dimensions: []*data.Dimension{
            {Name: "eventName"},
        },
        Metrics: []*data.Metric{
            {Name: "eventCount"},
            {Name: "totalUsers"},
        },
        DateRanges: []*data.DateRange{
            {
                StartDate: "2023-10-30",
                EndDate: "2023-11-26",
            },
        },

        DimensionFilter: &data.FilterExpression{
            Expr: &data.FilterExpression_Filter{
                Filter: dimensionFilter,
            },
        },
    }





	//Multiple value dimension filter
	// Create two filter expressions for each event name
    eventNameFilter1 := &data.Filter{
        FieldName: "eventName",
        OneFilter: &data.Filter_StringFilter_{
            StringFilter: &data.Filter_StringFilter{
                MatchType:     data.Filter_StringFilter_EXACT,
                Value:         "banner_click",
                CaseSensitive: false,
            },
        },
    }

    eventNameFilter2 := &data.Filter{
        FieldName: "eventName",
        OneFilter: &data.Filter_StringFilter_{
            StringFilter: &data.Filter_StringFilter{
                MatchType:     data.Filter_StringFilter_EXACT,
                Value:         "header_click",
                CaseSensitive: false,
            },
        },
    }

    // Create an OR group to combine the two filter expressions
    orGroup := &data.FilterExpression_OrGroup{
        OrGroup: &data.FilterExpressionList{
            Expressions: []*data.FilterExpression{
                {Expr: &data.FilterExpression_Filter{Filter: eventNameFilter1}},
                {Expr: &data.FilterExpression_Filter{Filter: eventNameFilter2}},
            },
        },
    }


    // Make a request to the GA4 Data API
    request := &data.RunReportRequest{
        Property: "properties/" + propertyID,
        Dimensions: []*data.Dimension{
            {Name: "eventName"},
        },
        Metrics: []*data.Metric{
            {Name: "eventCount"},
            {Name: "totalUsers"},
        },
        DateRanges: []*data.DateRange{
            {
                StartDate: "2023-10-30",
                EndDate: "2023-11-26",
            },
        },
        DimensionFilter: &data.FilterExpression{
            Expr: orGroup,
        },
    }




*/
