package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	pb "google.golang.org/genproto/googleapis/analytics/data/v1beta"
)

// Report represents a report structure
type Report struct {
	Dimensions []string `json:"dimensions"`
	Metrics    []string `json:"metrics"`
}

// AnalyticsClient defines the interface for interacting with Analytics Data API
type AnalyticsClient interface {
	RunReport(ctx context.Context, req *pb.RunReportRequest) (*pb.RunReportResponse, error)
}

// ServiceAccount struct holds the service account key information
type ServiceAccount struct {
	Type       string `json:"type"`
	ProjectID  string `json:"project_id"`
	PrivateKey string `json:"private_key"`
}

func NewClient(ctx context.Context, serviceAccountPath string) (AnalyticsClient, error) {
	// Open the service account key file
	f, err := os.Open(serviceAccountPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open service account key file: %w", err)
	}
	defer f.Close() // Ensure file is closed

	// Parse the service account key data
	var sa ServiceAccount
	if err := json.NewDecoder(f).Decode(&sa); err != nil {
		return nil, fmt.Errorf("failed to parse service account key: %w", err)
	}

	// Create credentials from the service account data
	cred, err := credentials.NewFromServiceAccountInfo(ctx, &sa, scopes...)
	if err != nil {
		return nil, fmt.Errorf("failed to create credentials: %w", err)
	}

	// Connect to gRPC server
	conn, err := grpc.DialContext(ctx, "analyticsdata.googleapis.com:443", grpc.WithTransportCredentials(cred))
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC server: %w", err)
	}

	return &analyticsClient{client: pb.NewBetaAnalyticsDataClient(conn)}, nil
}

// analyticsClient implements AnalyticsClient interface
type analyticsClient struct {
	client pb.BetaAnalyticsDataClient
}

// RunReport calls the RunReport API method
func (c *analyticsClient) RunReport(ctx context.Context, req *pb.RunReportRequest) (*pb.RunReportResponse, error) {
	return c.client.RunReport(ctx, req)
}

// RunReportWithCredentials calls RunReport with authentication using service account key
func RunReportWithCredentials(ctx context.Context, serviceAccountPath string, propertyID string, report *Report) (*pb.RunReportResponse, error) {
	// Create an Analytics client
	client, err := NewClient(ctx, serviceAccountPath)
	if err != nil {
		return nil, err
	}

	// Create RunReport request
	req := &pb.RunReportRequest{
		Property: fmt.Sprintf("properties/%s", propertyID),
		DateRanges: []*pb.DateRange{
			{StartDate: "7DaysAgo", EndDate: "today"},
		},
		Dimensions: report.Dimensions,
		Metrics:    report.Metrics,
	}

	// Add authorization header with service account email
	md := metadata.New(map[string][]string{"authorization": {"Bearer <your-service-account-email>"}})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Run the report
	return client.RunReport(ctx, req)
}
