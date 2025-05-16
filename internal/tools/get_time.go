package tools

import (
	"time"
)

type GetTimeInput struct{}

func init() {
	Register(ToolDefinition{
		Name: "get_time",
		Description: "Retrieve the current time as a formatted string " +
			"in ISO 8601 format (e.g., 2023-10-05T14:30:00Z)",
		Parameters: GenerateSchema[GetTimeInput](),
		Function: WithDecodedInput(func(val GetTimeInput) (string, error) {
			return time.Now().Format(time.RFC3339), nil
		}),
	})
}
