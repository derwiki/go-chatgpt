# OpenAI ChatGPT Command Line Interface

A Go command-line interface to communicate with OpenAI's ChatGPT API.
This program sends a prompt or question to the ChatGPT API for several models, prints the generated response for each,
and then sends all the responses to chatgpt-3.5-turbo to ask which is best.

## Example

```shell
➜  go-chatgpt git:(main) ✗ ./chatgpt "What is the meaning of life?"

> Chat Completion (gpt-3.5-turbo):
As an AI language model, I don't have a personal belief on the meaning of life. However, the meaning of life is a philosophical question that has baffled humans for centuries. It is a subjective concept that varies from person to person. Some may believe that the meaning of life is to pursue happiness and pleasure, while others may believe it to be a spiritual or religious pursuit. Ultimately, the meaning of life is personal and can only be determined by the individual.

> Chat Completion (text-davinci-003):
As an AI language model, I don't have personal beliefs and values, so I cannot hold my own meaning of life. The meaning of life can vary from person to person, and it's a philosophical and existential question that has been asked for centuries by philosophers, scientists, and spiritual leaders. Some believe it is to pursue happiness and personal fulfillment, others believe it is to fulfill a God-given purpose, while others see it as simply to exist and experience the universe around us.

> Chat Completion (text-davinci-002):
As an AI language model, I am incapable of having personal beliefs or opinions. However, the meaning of life is a highly debated philosophical and existential question that has no single, definitive answer. Different people may have different beliefs, interpretations, and perspectives regarding the purpose, significance, and value of life. Some may find meaning in spirituality, religion, or personal goals, while others may focus on human connection, experiences, or the pursuit of happiness. Ultimately, the meaning of life is a subjective and complex concept that may vary for each individual.

> Text Completion (da-vinci-002):
This is an individual question that each person has to answer for themselves. There is no one-size-fits-all answer to this question. Everyone has their own perspectives and interpretations, and it is these individual perspectives that help give our lives meaning. Ultimately, it is up to you to decide what the meaning of life is for you.

> Which of those answers is best?
All three answers provide good perspectives on the meaning of life, but the last one is the best because it emphasizes the importance of individual perspective and interpretation. It also highlights that there is no one definitive answer to this question and that it's up to each person to decide for themselves what the meaning of life is.
```

## Build
ChatGPT is built in Go, and requires Go 1.16+ to be installed. To build ChatGPT, run the following command:
```
go build -o chatgpt
```

## Usage
The ChatGPT CLI tool can accept a prompt either as a command-line argument or as standard input.

### Command-Line Argument
The following command sends a prompt as a command-line argument:
```shell
./chatgpt "What is the meaning of life?"
```

### Standard Input
The following command sends a prompt as standard input:
```shell
echo "What is the meaning of life?" | ./chatgpt
```

### Environment Variables
The following environment variables can be used to configure ChatGPT:

#### `OPENAI_API_KEY`
This variable is used to authenticate your OpenAI API key. If this variable is not set, ChatGPT will look for your key in the `.openai_key` file.

#### `MAX_TOKENS`
Defines the maximum number of tokens to generate in the response. The default value is `100`.

#### `PROMPT_PREFIX`

Defines a prefix that is prepended to any prompts sent to the API. Mostly useful when data is coming in on STDIN and you
want to add instructions preceding, e.g.:
```shell
➜  go-chatgpt git:(main) ✗ cat chatgpt.go | PROMPT_PREFIX="Suggest improvements for this Go program: " ./chatgpt

> Chat Completion (gpt-3.5-turbo):
Here are some potential improvements for this Go program:

1. Add error handling: There are several places in this program where errors can occur, such as when making API requests or parsing configuration values. It would be helpful to include more robust error handling to provide more helpful error messages and prevent the program from crashing.

2. Use a package manager: Instead of managing dependencies manually, it would be better to use a package manager like Go Modules to manage dependencies automatically.

3. Consolidate API response handling: The program currently has two separate functions for handling API requests and responses, which can make it harder to read and maintain. It might make sense to consolidate these functions into a single function that can handle requests and responses for any OpenAI API.

4. Create an API client package: Instead of having API functions scattered throughout the codebase, it would be better to create a separate package for OpenAI API clients that can be imported and reused across different programs.

5. Use environment variables for model selection: Rather than hardcoding the model names in the code, it would be better to use environment variables to specify which models to use. This would make the code more flexible and easier to configure.

6. Improve channel handling: The program currently launches multiple goroutines to call the API functions in parallel and waits for their results with channels. However, this approach can be difficult to manage and scale as the number of API requests grows. It might make more sense to use a worker pool or other concurrency pattern to handle API requests more efficiently.

> Chat Completion (text-davinci-003):
1. Remove unused imports, such as "strconv".

2. Use constants or variables for the model name in ChatGPTCompletionsRequest instead of hardcoding it.

3. Use context.Background() in getTextCompletion to be consistent with getChatCompletions.

4. Print the error message instead of calling log.Fatal in the getTextCompletion function.

5. Add error handling when loading the environment variables. Instead of returning an empty Config and error, it is better to return a Config with default values and the error message.

6. Remove unnecessary spaces in print statements.

7. Add comments to explain what each function does.

8. Use a switch statement instead of if else if for the model input to be more readable.

9. Improve the error message for when no prompt is found.

10. Use defer to close the response body in getChatCompletions instead of at the end of the function.

11. Use "strconv.Itoa" instead of fmt.Sprintf for the MaxTokens value in ChatGPTCompletionsRequest.

12. Rename variables with shorter and more descriptive names.

13. Remove the printUsage function since it is not used.

14. Add error handling when decoding the json response in getTextCompletion.

> Chat Completion (text-davinci-002):
1. Improve error handling: Right now, the program uses `log.Fatal` to handle errors, which terminates the program. It would be more user-friendly to return error messages and let the user decide how to handle them.

2. Make constants configurable: The program currently specifies the GPT models and API URL as constants. It would be better to make these configurable by the user, either via command line arguments or environment variables.

3. Simplify main function: The main function currently launches goroutines to call the API functions in parallel, which can be complex to understand. It would be simpler to use a loop to call the API functions sequentially and store the results in an array.

4. Refactor API request functions: The current implementation of the `getTextCompletion` function requires the caller to pass in a large `Config` struct. It would be simpler to pass in only the necessary parameters. Similarly, the `getChatCompletions` function could be simplified by allowing the caller to pass in the GPT model as a parameter, rather than requiring it to be hard-coded.

5. Improve user interface: The program currently prints the raw API responses to the console, which may be difficult for users to read. It would be better to format the responses and provide options for the user to refine or choose the best response.

6. Unit tests: The current program lacks unit tests, which can help catch bugs and improve code quality. It would be good to write unit tests for each of the API request functions and any helper functions.

> Text Completion (da-vinci-002):
Suggestions:
1. Separate the code into multiple functions, such as, getTextCompletion(), parseResponse(), printSummary() and such, to make the code more structurally organized and easier to read.
2. Add documentation to the code to explain how each function is used.
3. Refactor the code to make it more efficient by removing redundant loops and variable declarations.
4. Handle potential errors in a better way, for example, by printing

> Which of those answers is best?
error messages to the console instead of terminating the program.
5. Use a logger instead of printing to the console directly.
6. Add support for different response formats, such as JSON and XML.
7. Implement caching to reduce the number of API requests made.
8. Add more control over the output, such as the number of results returned or sorting options.
9. Use a configuration file to set program options, such as API credentials and model selection.
10. Validate user input to prevent errors and improve security.
```

## Contributing
If you want to contribute to ChatGPT, you can send a pull request with your changes. Before doing so, please make sure all tests pass by running the following command:
```
go test
```

## License
This project is released under the MIT License.