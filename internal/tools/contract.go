package tools

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
	"github.com/ollama/ollama/api"
)

type ToolDefinition struct {
	Name        string                                      `json:"name"`
	Description string                                      `json:"description"`
	Parameters  *jsonschema.Schema                          `json:"parameters"`
	Function    func(toolCall api.ToolCall) (string, error) `json:"-"`
}

func (def ToolDefinition) ToOllama() (api.Tool, error) {
	defJSON, err := json.Marshal(def)
	if err != nil {
		return api.Tool{}, err
	}

	var tf api.ToolFunction
	if err := json.Unmarshal(defJSON, &tf); err != nil {
		return api.Tool{}, err
	}

	return api.Tool{
		Type:     "function",
		Function: tf,
	}, nil
}

func GenerateSchema[T any]() *jsonschema.Schema {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T

	schema := reflector.Reflect(v)

	return schema
}

func WithDecodedInput[T any](fn func(val T) (string, error)) func(toolCall api.ToolCall) (string, error) {
	return func(toolCall api.ToolCall) (string, error) {
		encoded, err := json.Marshal(toolCall.Function.Arguments)
		if err != nil {
			return "", err
		}
		var val T
		if err := json.Unmarshal(encoded, &val); err != nil {
			return "", err
		}
		return fn(val)
	}
}
