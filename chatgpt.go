package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const apiBaseURL = "https://api.openai.com/v1/completions"

type ChatGPTCompletionsResponse struct {
	Choices []ChatGPTCompletionsResponseChoice `json:"choices"`
}
type ChatGPTCompletionsResponseChoice struct {
	FinishReason string `json:"finish_reason"`
	Index        int    `json:"index"`
	LogProbs     string `json:"logprobs"`
	Text         string `json:"text"`
}
type ChatGPTCompletionsRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type Config struct {
	OpenAIApiKey string
	MaxTokens    int
}

func getChatGPTResponse(prompt string, config Config) string {
	chatGPTCompletionsRequest := ChatGPTCompletionsRequest{
		Model:     "text-davinci-003",
		Prompt:    prompt,
		MaxTokens: config.MaxTokens,
	}
	requestBodyBytes, err := json.Marshal(chatGPTCompletionsRequest)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	request, err := http.NewRequest("POST", apiBaseURL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.OpenAIApiKey))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	// close the response body at the end of the function
	defer response.Body.Close()

	var responseBody ChatGPTCompletionsResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		log.Fatal(err)
	}

	if len(responseBody.Choices) == 0 {
		log.Fatal("No choices found in the response body.")
	}

	return strings.TrimSpace(responseBody.Choices[0].Text)
}

func getChatCompletions(content string, config Config) string {
	client := openai.NewClient(config.OpenAIApiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content
}

func hasStdinInput() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	return info.Mode()&os.ModeCharDevice == 0
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Fatal error occurred in loadConfig")
	} else if len(os.Args) > 1 {
		fmt.Println("> Using prompt from args:", os.Args[1])
		fmt.Println(getChatGPTResponse(os.Args[1], config))
	} else if hasStdinInput() {
		scanner := bufio.NewScanner(os.Stdin)

		scanner.Split(bufio.ScanBytes)
		var buffer bytes.Buffer
		for scanner.Scan() {
			buffer.Write(scanner.Bytes())
		}

		prompt := strings.TrimSpace(buffer.String())
		fmt.Println("> Using prompt from STDIN:", prompt)
		fmt.Println("> Chat Completion (gpt-3.5-turbo):", prompt)
		fmt.Println(getChatCompletions(prompt, config))
		fmt.Println("> Text Completion (da-vinci-002):", prompt)
		fmt.Println(getChatGPTResponse(prompt, config))
	} else {
		fmt.Println("X No prompt found in args or STDIN")
		printUsage()
	}
}

func loadConfig() (Config, error) {
	config := Config{}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKeyBytes, err := ioutil.ReadFile("./.openai_key")
		if err != nil {
			return config, err
		}
		apiKey = strings.TrimSpace(string(apiKeyBytes))
	}
	config.OpenAIApiKey = apiKey

	maxTokensStr := os.Getenv("MAX_TOKENS")
	if maxTokensStr == "" {
		config.MaxTokens = 100
	} else {
		maxTokens, err := strconv.Atoi(maxTokensStr)
		if err != nil {
			return config, err
		}
		config.MaxTokens = maxTokens
	}

	return config, nil
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
