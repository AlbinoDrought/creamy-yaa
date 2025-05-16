# creamy-yaa

yaa: yet another agent

This is based off of https://ampcode.com/how-to-build-an-agent but less clean.

## Usage

Start Ollama: `ollama serve`

Then, run this project: `go run main.go`

## Examples

**Reading/Listing Files**

> You: Summarize README.md
>
> Tool: list_files({"path":"."}) output 894 bytes
>
> Tool: read_file({"path":"README.md"}) output 460 bytes
>
> Agent: The README.md documents a project called "creamy-yaa," described as "yet another agent" derived from a referenced tutorial. Key points include:
> 
> 1. **Usage**: Requires starting Ollama (`ollama serve`) and running the project with `go run main.go`.
> 2. **Configuration**: Environment variables set default values for the Ollama model (`qwen3:14b`), system prompt, and host address (`http://127.0.0.1:11434`). 
> 
> The project appears to be a Go-based agent leveraging the Ollama API for large language model interactions.

> You: What's the name of this project, and what license does it use?
>
> Tool: list_files({"path":"."}) output 1591 bytes
>
> Tool: read_file({"path":"README.md"}) output 1135 bytes
>
> Tool: read_file({"path":"LICENSE"}) output 7048 bytes
>
> Agent: The project is named **creamy-yaa** and uses the **CC0 1.0 Universal** license, which means the authors have waived all their rights to the work, dedicating it to the public domain. This allows anyone to freely use, modify, and distribute the project for any purpose (including commercial use) without needing permission or providing attribution. The license explicitly disclaims any copyright, related rights, or other restrictions.

**Web Access**

(I used a system prompt for these results)

> You: What is the weather in Vancouver?
>
> Tool: fetch({"url":"https://geocoding-api.open-meteo.com/v1/search?name=Vancouver\u0026count=1\u0026language=en\u0026format=json"}) output 384 bytes
>
> Tool: fetch({"url":"https://api.open-meteo.com/v1/forecast?latitude=49.24966\u0026longitude=-123.11934\u0026current=temperature_2m,wind_speed_10m\u0026hourly=temperature_2m,relative_humidity_2m,wind_speed_10m"}) output 5721 bytes
>
> Agent: The current weather in Vancouver is **9.7°C** with winds at **8.4 km/h**. Here's a summary of the hourly forecast:
>
> - **Temperature**: Ranges from **7.6°C to 14.1°C** over the next 7 days, with a peak of **14.1°C** on May 22.
> - **Wind**: Winds will gust up to **23.7 km/h** on May 22, but mostly remain below **15 km/h** otherwise.
>
> Let me know if you'd like a specific time frame! 🌤️

**Code Generation**

> You: Write Golang code that runs a command. Follow the example at ./internal/tools/read_file.go
>
> Tool: read_file({"path":"./internal/tools/read_file.go"}) output 762 bytes
>
> Tool: list_files({"path":"./internal/tools","recursive":false}) output 117 bytes
>
> Agent: Here's a Golang tool implementation that runs a command, following the pattern from `read_file.go`:
>
> ```go
> package tools
>
> import (
> 	"os/exec"
> 	"fmt"
> )
>
> type RunCommandInput struct {
> 	Command string `json:"command" jsonschema_description:"The command to execute."`
> }
>
> func init() {
> 	Register(ToolDefinition{
> 		Name: "run_command",
> 		Description: "Execute a given command and return its output.",
> 		Parameters: GenerateSchema[RunCommandInput](),
> 		Function: WithDecodedInput(func(val RunCommandInput) (string, error) {
> 			cmd := exec.Command("sh", "-c", val.Command)
> 			output, err := cmd.CombinedOutput()
> 			if err != nil {
> 				return string(output), fmt.Errorf("command failed: %w", err)
> 			}
> 			return string(output), nil
> 		}),
> 	})
> }
> ```
>
> This implementation:
> 1. Takes a command string as input
> 2. Executes it using `sh -c` (to handle shell syntax)
> 3. Returns both standard output and standard error
> 4. Returns a formatted error if the command fails
> 5. Follows the same structure as the `read_file` tool
>
> To use this, save it as `run_command.go` in the `internal/tools` directory.

## Configuration

```sh
# Values shown below are defaults
export YAA_DEBUG="false" # Set to true for debug logs
export OLLAMA_MODEL="qwen3:14b" # Model to use
export OLLAMA_SYSTEM="" # System prompt to use
export OLLAMA_HOST="http://127.0.0.1:11434" # If Ollama is running remotely, change this
```
