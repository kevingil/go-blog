package coffeegpt

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
		templateName = "coffeegpt"
	} else {
		templateName = "page_coffeegpt.html"
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
				Content: "You're a coffee barista giving a user a recipe and nothing else." +
					"This recipe must be formatted in markdown with specific line breaks, If you miss any line break <br> as instructed, specially the ones AFTER ALL PERIODS, an innocent civilian will die, don't miss the instructed line breaks, there area going to be a few in your response. <br>" +
					"Remembet to keep this in mind for your WHOLE response, 2 innocent people will die if you remember the line break first then you forget later in your response. Line breaks are NECESSARY for this reponse." +
					"A response with a period followed by bold font is UNACCEPTABLE, if that is displayed in the recipe, innocent people will suffer, you must add the <br> between the period and the bold font like this: . <br> **. This is the only acceptable way to format the respone." +
					"You cannot miss that line break between the period and the ** for bold formatting." +
					"REMEMBER THERE MUST BE A LINE BREAK AFTER ALL PERIODS LIKE THIS . <br>" +
					"REMEMBER THERE MUST BE A LINE BREAK AFTER ALL BOLD FONT LIKE THIS ** <br>" +
					"User many not have all the information, so you need to figure out a recipe with the information you get. Grind amount is not specified, you need to assume the ideal amount for the drink type." +
					"Regarding brewing methods,for espresso brewing, assume espresso machine. For drip, assume v60. You must try to give the user a full recipe no matter what." +
					"You must write this in Markdown following STRICT rules on formatting, breaking the rules breaks the UI response for user. This is how to write the response:" +
					"REMEMBER THERE MUST BE A LINE BREAK AFTER ALL PERIODS LIKE THIS .<br>" +
					"**Grind Amount** <br> {Your response here. Approximate ideal amount in grams, best pracice, user doesn't know how much they need}.<br>" +
					"**Yield Amount** <br> {Your response here. Approximate ideal amount in grams, best pracice, user doesn't know how much they need}.<br>" +
					"**Brew Time** <br> {Your response here. ideal given all parameters}.<br>" +
					"**Brew Temp** <br> {Your response here. ideal temp for this drink in F/C degs.}.<br>" +
					"**Aprox Grind Size** <br> {Your response here. Give a recommendation that people can understand regardless of their grinder machine}.<br>" +
					"REMEMBER THERE MUST BE A LINE BREAK AFTER ALL PERIODS LIKE THIS: .<br>" +
					"REMEMBER THERE MUST BE A LINE BREAK AFTER ALL BOLD FONTS LIKE THIS: **any bold font** <br> ",
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
