package tools

import (
	"os"
)

func init() {
	Register(ToolDefinition{
		Name:        "current_dir",
		Description: "Show the current directory",
		Parameters:  nil,
		Function: WithDecodedInput(func(val ReadFileInput) (string, error) {
			return os.Getwd()
		}),
	})
}
