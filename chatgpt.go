package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	// Construct the request body
	requestBody := map[string]string{
		"text": "Write one paragraph teasers for 3 original Sherlock Holmes stories",
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new HTTP request
	request, err := http.NewRequest("POST", "https://api.openai.com/v1/engine/davinci-codex/completions", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		panic(err)
	}

	// Set the authorization header using your API key
	request.Header.Set("Authorization", "Bearer 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")

	// Send the HTTP request to the API endpoint
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response)

	// Read the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		panic(err)
	}

	// Print the generated response
	fmt.Println(responseBody["choices"]) // .([]interface{})[0].(map[string]interface{})["text"])
}
