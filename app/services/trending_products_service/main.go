package trending_products_service

/*
Documentation
https://developers.google.com/analytics/devguides/reporting/data/v1/quickstart-client-libraries#go
*/

import (
	"os"
)

func main() {

	//Service Variables
	URL_PROPERTY_ID := os.Getenv("URL_PROPERTY_ID")
	GA_API_KEY := os.Getenv("GA_API_KEY")

	property := property{URL_PROPERTY_ID, GA_API_KEY}

	// TO DO
	// FEED THIS TO AN AI

	property.getPopularProductsByViews()
	property.getPopularProductsByPurchaseAmount()

}
