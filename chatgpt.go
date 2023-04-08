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

func printUsage() {
	fmt.Println(`Usage:
  ./chatgpt [PROMPT]
  echo "PROMPT" | ./chatgpt

Description:
  A Golang client for OpenAI's ChatGPT API. This program takes a user prompt
  as a quoted command-line argument or via the standard input (STDIN), sends
  it to the API, and prints the generated response.

Options:
  PROMPT              The question or prompt to send to the ChatGPT API.

Environment Variables:
  OPENAI_API_KEY      Your OpenAI API key.
  MAX_TOKENS          The maximum number of tokens to generate in the response. (default: 100)

Example:
  ./chatgpt "What is the capital of France?"
  echo "What is the capital of France?" | ./chatgpt`)
}

func gptApiCall(prompt string, config Config) string {
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

	// Return the generated response
	return strings.TrimSpace(responseBody.Choices[0].Text)
}

func hasStdinInput() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	return info.Mode()&os.ModeCharDevice == 0
}

func main() {
	config := loadConfig()
	if len(os.Args) > 1 {
		fmt.Println("> Using prompt from args:", os.Args[1])
		fmt.Println(gptApiCall(os.Args[1], config))
	} else if hasStdinInput() {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		prompt := scanner.Text()
		fmt.Println("> Using prompt from STDIN:", prompt)
		fmt.Println(gptApiCall(prompt, config))
	} else {
		fmt.Println("X No prompt found in args or STDIN")
		printUsage()
	}
}
