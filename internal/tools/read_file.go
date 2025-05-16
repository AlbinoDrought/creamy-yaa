package tools

import "os"

type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

func init() {
	Register(ToolDefinition{
		Name:        "read_file",
		Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
		Parameters:  GenerateSchema[ReadFileInput](),
		Function: WithDecodedInput(func(val ReadFileInput) (string, error) {
			data, err := os.ReadFile(val.Path)
			return string(data), err
		}),
	})
}
