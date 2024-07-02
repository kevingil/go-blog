package controllers

import (
	"context"
	"fmt"
	"log"

	data "google.golang.org/genproto/googleapis/analytics/data/v1beta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var AnalyticsPropertyID string
var AnalyticsServiceAccountJsonPath string

func getAnalyticsData(req *data.RunReportRequest) (*data.RunReportResponse, error) {
	// Load service account credentials from JSON key file
	creds, err := oauth.NewServiceAccountFromFile(AnalyticsServiceAccountJsonPath, "https://www.googleapis.com/auth/analytics.readonly")
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

	fmt.Println("Report result:")
	for _, row := range response.Rows {
		// fmt.Printf("%s, %v\n", row.DimensionValues[0].GetValue(), row.MetricValues[0].GetValue())
		fmt.Printf("%s, Event Count: %v, Active Users: %v\n", row.DimensionValues[0].GetValue(), row.MetricValues[0].GetValue(), row.MetricValues[1].GetValue())
	}
	return response, nil
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
