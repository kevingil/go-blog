package espresso

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kevingil/blog/app/templates"
)

/*
func init() {

	 Init router
	when running standalone app
	router.Init()


	// Init DB connection
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	maxRetries := 5
	retryInterval := 5 * time.Second
	for i := 0; i < maxRetries; i++ {
		//models.Db, models.Err = sql.Open("mysql", dsn)
		models.Db, models.Err = sql.Open("mysql", os.Getenv("ESPRESSO_APP_NEON"))

		if models.Err == nil {
			break
		}

		fmt.Printf("Failed to connect to MySQL server: %v\n", models.Err)
		fmt.Printf("Retrying in %v...\n", retryInterval)
		time.Sleep(retryInterval)
	}
}
*/

// Web entry point
func EspressoApp(w http.ResponseWriter, r *http.Request) {
	isHTMXRequest := r.Header.Get("HX-Request") == "true"
	var templateName string

	if isHTMXRequest {
		templateName = "espressoai"
	} else {
		templateName = "page_espressoai.html"
	}

	var response bytes.Buffer

	if err := templates.Tmpl.ExecuteTemplate(&response, templateName, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())
}

// Your OpenAI API key
const apiKey = "YOUR_API_KEY"

func GenerateResponse(w http.ResponseWriter, r *http.Request) {
	// Read the question from the request body
	questionBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	question := string(questionBytes)

	// Define the OpenAI API endpoint
	apiUrl := "https://api.openai.com/v1/engines/davinci-codex/completions"

	// Create a POST request to the OpenAI API
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer([]byte(fmt.Sprintf(`{
        "prompt": "%s",
        "max_tokens": 50
    }`, question))))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Set the API key in the request header
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request to OpenAI
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request to OpenAI", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read and send the response back to the client
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from OpenAI", http.StatusInternalServerError)
		return
	}

	// Send the OpenAI response as the HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBytes)
}

func LoadAIResponse(w http.ResponseWriter, r *http.Request) {

}
