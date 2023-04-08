# ChatGPT Golang Client

This is a simple Golang client for interacting with OpenAI's ChatGPT API. It takes a user prompt from the standard input (STDIN), sends it to the API, and prints the generated response.
Prerequisites

 * Golang installed (version 1.16+ is recommended)
 * OpenAI API key

## Installation

Clone this repository:

    git clone https://github.com/derwiki/go-chatgpt.git

Navigate to the project directory:

    cd go-chatgpt

Create a .openai_key file in the project directory and paste your OpenAI API key inside:

    echo "your_openai_api_key" > ./.openai_key

## Usage

Build the program:

    go build chatgpt.go

Run the compiled binary, providing the user prompt as input:


    echo "Your prompt here" | ./chatgpt-client

You can also change the number of generated tokens by setting the MAX_TOKENS environment variable:

    MAX_TOKENS=50 echo "Your prompt here" | ./chatgpt-client

The generated response from ChatGPT will be printed to the standard output.

## License

This project is released under the MIT License.