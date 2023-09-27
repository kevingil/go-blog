package trending_products_service

func ClientStart() {
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
}
