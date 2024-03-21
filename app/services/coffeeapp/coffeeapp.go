package coffeeapp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/kevingil/blog/app/controllers"
	"github.com/kevingil/blog/app/views"

	openai "github.com/sashabaranov/go-openai"
)

// Render template
func CoffeeApp(w http.ResponseWriter, r *http.Request) {
	views.Hx(w, r, "main_layout", "coffeeapp", controllers.Context{})

}

// Render completion
func Completion(w http.ResponseWriter, r *http.Request) {
	isHTMXRequest := r.Header.Get("HX-Request") == "true"

	var response bytes.Buffer

	if isHTMXRequest {
		if err := views.Tmpl.ExecuteTemplate(&response, "coffeeapp_completion", nil); err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())

}

// Stream Chat Completion
func StreamRecipe(w http.ResponseWriter, r *http.Request) {
	// User prompt
	question := r.URL.Query().Get("question")

	// OpenAI client.
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Create request
	request := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 300,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "You're a coffee barista tasked with providing a coffee bean recipe in markdown format. The line breaks <br> are important to format the response." +
					"The user will specify the coffee type, bean process, growing elevation, color, and bean origin in their message." +
					"Your response should follow this html format, keep responses short, just the measurements:" +
					"<p><b>Grind Amount</b>[Provide an approximate coffee bean grind amount in grams, and I'd like it to be accurate within +/-1 gram.].</p>" +
					"<p><b>Yield Amount</b>[Provide an approximate extraction yield amount in grams, and I'd like it to be accurate within +/-1 gram.].</p>" +
					"<p><b>Brew Time</b>[How long the brew time should take within +/-3sec given all parameters].</p>" +
					"<p><b>Brew Temp</b>[Ideal brew temperature given all parameters in F/C degrees].</p>" +
					"Some tips to keep in mind: higher grown beans tend to be more dense, darker roasts result in beans to be more brittle / less dense, the more intense the processing method, the more brittle the bean becomes. Density for common processing methods: Washed > Semiwashed > Honey > Natural",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question,
			},
		},
		Stream: true,
	}

	stream, err := client.CreateChatCompletionStream(context.Background(), request)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	// Start SSE event
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	fmt.Printf("event: coffee-help\n")
	fmt.Fprintf(w, "event: coffee-help\n")

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Printf("\nEnd of stream.\n")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		// Send SSE completion data and log for debugging
		fmt.Printf("data: %s\n\n", response.Choices[0].Delta.Content)
		fmt.Fprintf(w, "data: %s\n\n", response.Choices[0].Delta.Content)
		w.(http.Flusher).Flush()
	}
}
