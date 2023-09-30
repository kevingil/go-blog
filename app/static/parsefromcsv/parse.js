
document.addEventListener("DOMContentLoaded", function () {
    const orderZipInput = document.getElementById("orderzip");
    const orderTaxRateSpan = document.getElementById("ordertaxrate");

    orderZipInput.addEventListener("change", function () {
        const zipCode = orderZipInput.value;

        // Load the CSV file directly from the same domain
        fetch("avalaratax.csv") // Adjust the path as needed
            .then(response => response.text())
            .then(csvData => {
                // Split the CSV data into rows
                const csvRows = csvData.split("\n");

                // Create an object to store tax rate data
                const taxRates = {};

                // Iterate through CSV rows and populate the taxRates object
                for (let i = 1; i < csvRows.length; i++) {
                    const row = csvRows[i].split(",");
                    if (row.length === 2) {
                        const zip = row[0].trim();
                        const rate = parseFloat(row[1].trim());

                        if (!isNaN(rate)) {
                            taxRates[zip] = rate;
                        }
                    }
                }

                // Check if the zip code exists in the taxRates object
                if (taxRates.hasOwnProperty(zipCode)) {
                    // Found the matching zip code, update tax rate
                    orderTaxRateSpan.textContent = taxRates[zipCode].toFixed(4);
                } else {
                    // If zip code not found, display a message
                    orderTaxRateSpan.textContent = "Tax TBD";
                }
            })
            .catch(error => {
                console.error("Error loading Tax data", error);
            });
    });
});
