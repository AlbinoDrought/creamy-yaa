package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/sirupsen/logrus"
	"go.albinodrought.com/creamy-yaa/internal/tools"
)

var logger = logrus.New()

func main() {
	if os.Getenv("YAA_DEBUG") == "true" {
		logger.SetLevel(logrus.DebugLevel)
		logger.Debug("debug mode")
	}

	ctx := context.Background()

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "qwen3:14b"
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		logger.WithError(err).Fatal("failed to api.ClientFromEnvironment")
	}

	apiTools, err := tools.LibraryToOllama()
	if err != nil {
		logger.WithError(err).Fatal("failed to create tools")
	}

	scanner := bufio.NewScanner(os.Stdin)

	baseMessages := []api.Message{} // we reset to this on !clear
	{
		system := os.Getenv("OLLAMA_SYSTEM")
		if system != "" {
			system = strings.ReplaceAll(system, `%YAA_DATE%`, time.Now().Format(time.RFC3339))
			baseMessages = append(baseMessages, api.Message{
				Role:    "system",
				Content: system,
			})
		}
	}
	messages := slices.Clone(baseMessages)
	name := "Agent"

	{
		fmt.Print("\u001b[92mYAA\u001b[0m: Getting the agent ready... ")
		for range 3 {
			resp, err := queryStruct[struct {
				Name string `json:"name"`
				UA   string `json:"ua"`
			}](ctx, client, model, messages, "Identify yourself, and also quickly generate an HTTP user-agent that includes your name.")

			if err != nil {
				logger.WithError(err).Debug("failed to request UA")
			}
			if resp.Name != "" && resp.UA != "" {
				name = resp.Name
				tools.FetchUserAgent = resp.UA
				logger.WithField("id", resp).Debug("set identity")
				break
			}
		}
		fmt.Print("\u001b[2K\r")
		fmt.Printf("\u001b[92mYAA\u001b[0m: Name: %v | Model: %v | UA: %v\n", name, model, tools.FetchUserAgent)
		fmt.Print("\u001b[92mYAA\u001b[0m: Reset history with !clear, retry last prompt with !retry\n")
		fmt.Println()
	}

	var output strings.Builder
	var skipInput bool
	var agentSpoke bool
	var agentFinishedThinking bool
	for {
		if !skipInput {
			fmt.Print("\u001b[94mYou\u001b[0m: ")
			if !scanner.Scan() {
				break
			}

			userContent := scanner.Text()
			if userContent == "!clear" {
				messages = slices.Clone(baseMessages)
				continue
			}
			if userContent == "!retry" {
				for i := len(messages) - 1; i >= 0; i-- {
					if messages[i].Role == "user" {
						messages = messages[:i+1]
						break
					}
				}
				skipInput = true
				continue
			}
			messages = append(messages, api.Message{
				Role:    "user",
				Content: userContent,
			})
		}
		skipInput = false

		output.Reset()
		agentSpoke = false
		agentFinishedThinking = false
		pendingToolCalls := []api.ToolCall{}
		err = client.Chat(ctx, &api.ChatRequest{
			Model:     model,
			Messages:  messages,
			Tools:     apiTools,
			KeepAlive: &api.Duration{Duration: 10 * time.Minute},
		}, func(cr api.ChatResponse) error {
			for _, toolCall := range cr.Message.ToolCalls {
				pendingToolCalls = append(pendingToolCalls, toolCall)
			}
			if cr.Message.Content != "" {
				output.WriteString(cr.Message.Content)
				if !agentFinishedThinking {
					endOfThoughtsIdx := strings.Index(cr.Message.Content, "</think>")
					if endOfThoughtsIdx != -1 {
						agentFinishedThinking = true
						cr.Message.Content = strings.TrimSpace(cr.Message.Content[endOfThoughtsIdx+len("</think>"):])
					}
				}

				if agentFinishedThinking {
					if !agentSpoke {
						fmt.Printf("\u001b[93m%v\u001b[0m: ", name)
						agentSpoke = true
					}
					fmt.Print(cr.Message.Content)
				}
			}
			return nil
		})
		if agentSpoke {
			fmt.Println()
		}
		if err != nil {
			logger.WithError(err).Error("failed to get response")
			continue
		}
		messages = append(messages, api.Message{
			Role:      "assistant",
			Content:   output.String(),
			ToolCalls: pendingToolCalls,
		})

		if len(pendingToolCalls) > 0 {
			skipInput = true
			for _, toolCall := range pendingToolCalls {
				fmt.Printf("\u001b[92mTool\u001b[0m: %v(%v) ", toolCall.Function.Name, toolCall.Function.Arguments.String())
				// messages = append(messages, api.Message{
				// 	Role:    "tool_call",
				// 	Content: fmt.Sprintf("%v(%v)", toolCall.Function.Name, toolCall.Function.Arguments),
				// })

				handler, ok := tools.Library[toolCall.Function.Name]
				if !ok {
					fmt.Print("\u001b[31merror\u001b[0m: not found\n")
					messages = append(messages, api.Message{
						Role:    "tool",
						Content: "error: tool not found",
					})
					continue
				}

				result, err := handler.Function(toolCall)
				if err != nil {
					fmt.Printf("\u001b[31merror\u001b[0m: %v\n", err)
					messages = append(messages, api.Message{
						Role:    "tool",
						Content: fmt.Sprintf("error: %v", err.Error()),
					})
					continue
				}

				fmt.Printf("output %v bytes\n", len(result))
				messages = append(messages, api.Message{
					Role:    "tool",
					Content: result,
				})
			}
		}
	}
}

func queryStruct[T any](
	ctx context.Context,
	client *api.Client,
	model string,
	messages []api.Message,
	request string,
) (T, error) {
	schema := tools.GenerateSchema[T]()
	schemaStr, err := json.Marshal(schema)
	if err != nil {
		return *new(T), err
	}

	var resp api.ChatResponse
	f := false
	client.Chat(ctx, &api.ChatRequest{
		Model: model,
		Messages: append(slices.Clone(messages), api.Message{
			Role:    "user",
			Content: fmt.Sprintf("%v. Output your response following the JSON schema %v", request, schemaStr),
		}),
		Stream:    &f,
		KeepAlive: &api.Duration{Duration: 10 * time.Minute},
	}, func(cr api.ChatResponse) error {
		resp = cr
		return nil
	})
	endOfThoughts := strings.LastIndex(resp.Message.Content, "</think>")
	if endOfThoughts != -1 {
		resp.Message.Content = resp.Message.Content[endOfThoughts+len("</think>"):]
	}
	resp.Message.Content = strings.TrimSpace(resp.Message.Content)

	var result T
	if err := json.Unmarshal([]byte(resp.Message.Content), &result); err != nil {
		return *new(T), err
	}
	return result, nil
}
