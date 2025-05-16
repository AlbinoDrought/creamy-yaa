package tools

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
)

type ListFileInput struct {
	Path      string `json:"path" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
	Recursive bool   `json:"recursive" jsonschema_description:"If true, also include files from all subdirectories"`
}

func init() {
	Register(ToolDefinition{
		Name:        "list_files",
		Description: "List files and directories at a given path. If no path is provided, lists files in the current directory.",
		Parameters:  GenerateSchema[ListFileInput](),
		Function: WithDecodedInput(func(val ListFileInput) (string, error) {
			if val.Path != "" && val.Path[0] == '~' {
				homeDir, err := os.UserHomeDir()
				if err == nil {
					val.Path = path.Join(homeDir, val.Path[1:])
				}
			}
			if val.Path == "" {
				val.Path = "."
			}

			var files []string
			if val.Recursive {
				err := filepath.Walk(val.Path, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					relPath, err := filepath.Rel(val.Path, path)
					if err != nil {
						return err
					}

					if relPath != "." {
						if info.IsDir() {
							files = append(files, relPath+"/")
						} else {
							files = append(files, relPath)
						}
					}
					return nil
				})
				if err != nil {
					return "", err
				}
			} else {
				res, err := os.ReadDir(val.Path)
				if err != nil {
					return "", err
				}
				for _, r := range res {
					if r.Type().IsDir() {
						files = append(files, r.Name()+"/")
					} else {
						files = append(files, r.Name())
					}
				}
			}

			result, err := json.Marshal(files)
			if err != nil {
				return "", err
			}

			return string(result), nil
		}),
	})
}
