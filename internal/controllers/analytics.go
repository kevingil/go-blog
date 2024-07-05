package controllers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
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
