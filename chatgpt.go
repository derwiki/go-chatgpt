package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	secretBytes, err := ioutil.ReadFile("./.openai_key")
	if err != nil {
		panic(err)
	}
	secret := strings.ReplaceAll(string(secretBytes), "\n", "")
	fmt.Println(string(secret))

	// Construct the request body
	requestBody := map[string]interface{}{
		"model":      "text-davinci-003",
		"prompt":     "Write a one paragraph teaser for an original Sherlock Holmes stories",
		"max_tokens": 100,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
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
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", secret))
	request.Header.Set("Content-Type", "application/json")

	// Send the HTTP request to the API endpoint
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	// fmt.Println(fmt.Sprintln("Response:\n%s", response))

	// Read the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		panic(err)
	}

	// Print the generated response
	fmt.Println(responseBody["choices"]) // .([]interface{})[0].(map[string]interface{})["text"])
	fmt.Println(responseBody["text"])    // .([]interface{})[0].(map[string]interface{})["text"])
}
