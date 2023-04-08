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
	"strings"
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

type Config struct {
	OpenAIApiKey string
	MaxTokens    int
}

func loadConfig() Config {
	config := Config{}

	// Read OpenAI API Key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKeyBytes, err := ioutil.ReadFile("./.openai_key")
		if err != nil {
			panic(err)
		}
		apiKey = strings.TrimSpace(string(apiKeyBytes))
	}
	config.OpenAIApiKey = apiKey

	// Read MaxTokens
	maxTokensStr := os.Getenv("MAX_TOKENS")
	if maxTokensStr == "" {
		config.MaxTokens = 100
	} else {
		maxTokens, err := strconv.Atoi(maxTokensStr)
		if err != nil {
			panic(err)
		}
		config.MaxTokens = maxTokens
	}

	return config
}

func main() {
	config := loadConfig()
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	prompt := scanner.Text()
	fmt.Println("Using ChatGPT prompt from STDIN:", prompt)

	// Construct the request body
	chatGPTCompletionsRequest := ChatGPTCompletionsRequest{
		Model:     "text-davinci-003",
		Prompt:    prompt,
		MaxTokens: config.MaxTokens,
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
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.OpenAIApiKey))
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
