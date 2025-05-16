package tools

import (
	"os"
	"path"
)

type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

func init() {
	Register(ToolDefinition{
		Name:        "read_file",
		Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
		Parameters:  GenerateSchema[ReadFileInput](),
		Function: WithDecodedInput(func(val ReadFileInput) (string, error) {
			if val.Path != "" && val.Path[0] == '~' {
				homeDir, err := os.UserHomeDir()
				if err == nil {
					val.Path = path.Join(homeDir, val.Path[1:])
				}
			}
			data, err := os.ReadFile(val.Path)
			return string(data), err
		}),
	})
}
