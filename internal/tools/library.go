package tools

import "github.com/ollama/ollama/api"

var Library = map[string]ToolDefinition{}

func Register(defs ...ToolDefinition) {
	for _, def := range defs {
		Library[def.Name] = def
	}
}

func LibraryToOllama() (api.Tools, error) {
	tools := make([]api.Tool, len(Library))
	var i int
	var err error
	for _, tool := range Library {
		tools[i], err = tool.ToOllama()
		i++
		if err != nil {
			return api.Tools{}, err
		}
	}
	return tools, nil
}
