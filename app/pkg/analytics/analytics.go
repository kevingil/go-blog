package analytics

import (
	"context"
	"log"

	ga "google.golang.org/genproto/googleapis/analytics/data/v1beta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

// Report represents a report structure
type Report struct {
	Dimensions []string `json:"dimensions"`
	Metrics    []string `json:"metrics"`
}

type Response struct {
	Data any // stil working on this
}

// AnalyticsClient defines the interface for interacting with Analytics Data API
type AnalyticsClient interface {
	RunReport(ctx context.Context, req *ga.RunReportRequest) (*ga.RunReportResponse, error)
}

// NewClient creates a new Analytics client with authentication
func NewClient(ctx context.Context, serviceAccountPath string) (AnalyticsClient, error) {

	// Load service account credentials from JSON key file
	creds, err := oauth.NewServiceAccountFromFile(serviceAccountPath, "https://www.googleapis.com/auth/analytics.readonly")
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
	client := ga.NewBetaAnalyticsDataClient(conn)

	// Wrap the client in your custom AnalyticsClient
	return &analyticsClient{client: client}, nil
}

// analyticsClient implements AnalyticsClient interface
type analyticsClient struct {
	client ga.BetaAnalyticsDataClient
}

// RunReport calls the RunReport API method
func (c *analyticsClient) RunReport(ctx context.Context, req *ga.RunReportRequest) (*ga.RunReportResponse, error) {
	return c.client.RunReport(ctx, req)
}

// GetReportData extracts desired data from the RunReportResponse (modify for your needs)
func GetReportData(response *ga.RunReportResponse) ([]map[string]interface{}, error) {
	if response.Rows == nil {
		return nil, nil // handle empty response
	}

	data := make([]map[string]interface{}, 0)
	for _, row := range response.Rows {
		rowMap := make(map[string]interface{})
		for _, dimension := range response.MetricHeaders {
			rowMap[dimension.Name] = row.GetDimensionValues()
		}
		data = append(data, rowMap)
	}
	return data, nil
}
