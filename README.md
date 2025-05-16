# creamy-yaa

yaa: yet another agent

This is based off of https://ampcode.com/how-to-build-an-agent but less clean.

## Usage

Start Ollama: `ollama serve`

Then, run this project: `go run main.go`

## Examples

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

## Configuration

```sh
# Values shown below are defaults
export OLLAMA_MODEL="qwen3:14b" # Model to use
export OLLAMA_SYSTEM="" # System prompt to use
export OLLAMA_HOST="http://127.0.0.1:11434" # If Ollama is running remotely, change this
```
