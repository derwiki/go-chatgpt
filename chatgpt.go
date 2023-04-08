package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type ChatGPTCompletionsResponse struct {
	Choices []ChatGPTCompletionsResponseChoice `json:"choices"`
}
type ChatGPTCompletionsResponseChoice struct {
	FinishReason string `json:"finish_reason""`
	Index        int    `json:"index""`
	LogProbs     string `json:"logprobs""`
	Text         string `json:"text""`
}
type ChatGPTCompletionsRequest struct {
	Model     string `json:"model""`
	Prompt    string `json:"prompt""`
	MaxTokens int    `json:"max_tokens""`
}

// LoadApiKey it's not reading this correctly
func LoadApiKey() string {
	maybeFromEnv := os.Getenv("OPENAI_API_KEY")
	if len(maybeFromEnv) > 0 {
		fmt.Println(fmt.Sprintf("using API key from env var: %s", maybeFromEnv))
		return maybeFromEnv
	}
	secretBytes, err := ioutil.ReadFile("./.openai_key")
	if err != nil {
		fmt.Println(fmt.Sprintf("using API key from ./.openai_key: %s", string(secretBytes)))
		return string(secretBytes)
	}

	panic("no api key")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	prompt := scanner.Text()
	fmt.Println("Using ChatGPT prompt from STDIN:", prompt)

	// Construct the request body
	// TODO(derwiki) use ChatGPTCompletionsRequest struct
	maxTokens, err := strconv.Atoi(os.Getenv("MAX_TOKENS"))
	if err != nil {
		maxTokens = 100
	}
	chatGPTCompletionsRequest := ChatGPTCompletionsRequest{
		Model:     "text-davinci-003",
		Prompt:    prompt,
		MaxTokens: maxTokens,
	}
	requestBodyBytes, err := json.Marshal(chatGPTCompletionsRequest)
	if err != nil {
		panic(err)
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new HTTP request
	request, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		panic(err)
	}

	// Set the authorization header using your API key
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", LoadApiKey()))
	request.Header.Set("Content-Type", "application/json")

	// Send the HTTP request to the API endpoint
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	// Read the response body
	var responseBody ChatGPTCompletionsResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		panic(err)
	}

	// Print the generated response
	fmt.Println(responseBody.Choices[0].Text)
}
