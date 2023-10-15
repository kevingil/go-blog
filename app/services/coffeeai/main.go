package coffeeai

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kevingil/blog/app/templates"
	openai "github.com/sashabaranov/go-openai"
)

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
}

func CoffeeApp(w http.ResponseWriter, r *http.Request) {
	isHTMXRequest := r.Header.Get("HX-Request") == "true"
	var templateName string

	if isHTMXRequest {
		templateName = "coffeeai"
	} else {
		templateName = "page_coffeeai.html"
	}

	var response bytes.Buffer

	if err := templates.Tmpl.ExecuteTemplate(&response, templateName, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, response.String())

}

func StreamRecipe(w http.ResponseWriter, r *http.Request) {
	// Get the question from the request URL.
	question := r.URL.Query().Get("question")

	// Create an OpenAI client.
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Create a chat completion request.
	request := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 150,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "You're a coffee barista giving a user a recipe and nothing else." +
					"REMEMBER THERE MUST BE A LINE BREAK AFTER ALL PERIODS LIKE THIS .<br>" +
					"REMEMBER THERE MUST BE A LINE BREAK AFTER ALL BOLD FONT LIKE THIS **<br>" +
					"User many not have all the information, so you need to figure out a recipe with the information you get. Grind amount is not specified, you need to assume the ideal amount for the drink type." +
					"Regarding brewing methods,for espresso brewing, assume espresso machine. For drip, assume v60. You must try to give the user a full recipe no matter what." +
					"You must write this in Markdown folling STRICT rules or the formatting will get messed up. This is how to write the response:" +
					"REMEMBER THERE MUST BE A LINE BREAK AFTER ALL PERIODS LIKE THIS .<br>" +
					"1. Grind Amount. Formatting: **Grind Amount** <br> {response}.<br>" +
					"2. Yield Amount. Formatting: **Yield Amount** <br> {response}.<br>" +
					"3. Brew Time. Formatting: **Brew Time** <br> {response}.<br>" +
					"3. Brew Temp. Formatting: **Brew Temp** <br> {response}.<br>" +
					"3. Aprox Grind Size. Formatting: **Aprox Grind Size** <br> {response}.<br>" +
					"REMEMBER THERE MUST BE A LINE BREAK AFTER ALL PERIODS LIKE THIS .<br>" +
					"REMEMBER THERE MUST BE A LINE BREAK AFTER ALL BOLD FONT LIKE THIS **<br>" +
					"Make sure you are keeping all {responses} shot and using <br> as I showed you or the recipe will break. No other response format is acceptable." +
					"Terrible things will happen if you forget that there is 2 spaces between 2 bold titles as such ** <br> <br> **" +
					"REMEMBER THERE MUST BE A LINE BREAK AFTER ALL PERIODS LIKE THIS .<br>" +
					"REMEMBER THERE MUST BE A LINE BREAK AFTER ALL BOLD FONT LIKE THIS **<br>",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question,
			},
		},
		Stream: true,
		Stop:   []string{"For a  coffee with  processed beans grown at an elevation of  and a bean color of rgb(, I would suggest the following:-"},
	}

	stream, err := client.CreateChatCompletionStream(context.Background(), request)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Send the SSE event with the response data.
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

		// Send stream completion via SSE
		fmt.Printf("data: %s\n\n", response.Choices[0].Delta.Content)
		fmt.Fprintf(w, "data: %s\n\n", response.Choices[0].Delta.Content)
		w.(http.Flusher).Flush()
	}
}
