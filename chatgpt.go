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

type TextCompletionResponse struct {
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
	PromptPrefix string
	Model        string
}

func getTextCompletion(prompt string, config Config) string {
	textCompletionRequest := ChatGPTCompletionsRequest{
		Model:     "text-davinci-003",
		Prompt:    config.PromptPrefix + prompt,
		MaxTokens: config.MaxTokens,
	}
	requestBodyBytes, err := json.Marshal(textCompletionRequest)
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

	var responseBody TextCompletionResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		log.Fatal(err)
	}

	if len(responseBody.Choices) == 0 {
		log.Fatal("No choices found in the response body.")
	}

	return strings.TrimSpace(responseBody.Choices[0].Text)
}

func getChatCompletions(content string, config Config, model string) string {
	if model == "" {
		model = openai.GPT3Dot5Turbo
	}
	// TODO(derwiki) assert model exists in openai package
	client := openai.NewClient(config.OpenAIApiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: config.PromptPrefix + content,
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

	var prompt string
	if err != nil {
		fmt.Println("error: Fatal occurred in loadConfig")
	} else if len(os.Args) > 1 {
		prompt = os.Args[1]
	} else if hasStdinInput() {
		scanner := bufio.NewScanner(os.Stdin)

		scanner.Split(bufio.ScanBytes)
		var buffer bytes.Buffer
		for scanner.Scan() {
			buffer.Write(scanner.Bytes())
		}

		prompt = strings.TrimSpace(buffer.String())
	} else {
		fmt.Println("error: No prompt found in args or STDIN")
		printUsage()
		return
	}

	// Create channels for the API responses
	gpt3TurboCh := make(chan string)
	gpt3Davinci003Ch := make(chan string)
	gpt3Davinci002Ch := make(chan string)
	textDavinci002Ch := make(chan string)

	// if a model is specified, only call that model and exit
	if config.Model != "" {
		if config.Model == openai.GPT3Dot5Turbo {
			fmt.Println(getChatCompletions(prompt, config, openai.GPT3Dot5Turbo))
		} else if config.Model == openai.GPT3TextDavinci003 {
			fmt.Println(getChatCompletions(prompt, config, openai.GPT3TextDavinci003))
		} else if config.Model == openai.GPT3TextDavinci002 {
			fmt.Println(getChatCompletions(prompt, config, openai.GPT3TextDavinci002))
		} else if config.Model == "text-davinci-002" {
			fmt.Println(getTextCompletion(prompt, config))
		}
		return
	}

	// Launch goroutines to call the API functions in parallel
	go func() {
		gpt3TurboCh <- getChatCompletions(prompt, config, openai.GPT3Dot5Turbo)
	}()
	go func() {
		gpt3Davinci003Ch <- getChatCompletions(prompt, config, openai.GPT3TextDavinci003)
	}()
	go func() {
		gpt3Davinci002Ch <- getChatCompletions(prompt, config, openai.GPT3TextDavinci002)
	}()
	go func() {
		textDavinci002Ch <- getTextCompletion(prompt, config)
	}()

	// Wait for the API responses from the channels
	gpt3TurboRes := <-gpt3TurboCh
	gpt3Davinci003Res := <-gpt3Davinci003Ch
	gpt3Davinci002Res := <-gpt3Davinci002Ch
	textDavinci002Res := <-textDavinci002Ch

	// TODO(derwiki) put this in config
	verbose := false
	if verbose {
		fmt.Println(prompt)
	}
	// Print the API responses
	fmt.Println("\n> Chat Completion (gpt-3.5-turbo):")
	fmt.Println(gpt3TurboRes)
	fmt.Println("\n> Chat Completion (text-davinci-003):")
	fmt.Println(gpt3Davinci003Res)
	fmt.Println("\n> Chat Completion (text-davinci-002):")
	fmt.Println(gpt3Davinci002Res)
	fmt.Println("\n> Text Completion (da-vinci-002):")
	fmt.Println(textDavinci002Res)

	refine := fmt.Sprintf("Which of the following answers is best? \n\n%s\n\n%s\n\n%s\n\n%s", gpt3TurboRes, gpt3Davinci003Res, gpt3Davinci002Res, textDavinci002Res)
	refined := getChatCompletions(refine, config, openai.GPT3Dot5Turbo)
	fmt.Println("\n> Which of those answers is best?")
	fmt.Println(refined)
}

func loadConfig() (Config, error) {
	config := Config{}

	config.PromptPrefix = os.Getenv("PROMPT_PREFIX")

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

	config.Model = os.Getenv("GPT_MODEL")

	return config, nil
}

func printUsage() {
	fmt.Println(`
Usage:
  ./chatgpt [PROMPT]
  echo "PROMPT" | ./chatgpt
  cat chatgpt.go | PROMPT_PREFIX="Improve this program" ./chatgpt

Description:
  A Go command-line interface to communicate with OpenAI's ChatGPT API.
  This program sends a prompt or question to the ChatGPT API for several models,
  prints the generated response for each, and then sends all the responses to
  chatgpt-3.5-turbo to ask which is best.

Required Options:
  PROMPT              The question or prompt to send to the ChatGPT API.

Environment Variables:
  OPENAI_API_KEY      Your OpenAI API key.
  MAX_TOKENS          The maximum number of tokens to generate in the response. (default: 100)
  PROMPT_PREFIX       A prefix to add to each prompt.
  GPT_MODEL           The model to use. If not specified, all models will be used.

Example:
  ./chatgpt "What is the capital of France?"

  > Chat Completion (gpt-3.5-turbo):
  The capital of France is Paris.

  > Chat Completion (text-davinci-003):
  The capital of France is Paris.

  > Chat Completion (text-davinci-002):
  The capital of France is Paris.

  > Text Completion (da-vinci-002):
  Paris.

  > Which of those answers is best?
  The first answer is the best as it includes a complete sentence and clear
  statement of the capital of France. The other answers are incomplete sentences
  or single words that do not provide enough information.
	`)
}
