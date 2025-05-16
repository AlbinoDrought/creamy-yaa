package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/sirupsen/logrus"
	"go.albinodrought.com/creamy-yaa/internal/tools"
)

var logger = logrus.New()

func main() {
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

	messages := []api.Message{}
	{
		system := os.Getenv("OLLAMA_SYSTEM")
		if system != "" {
			messages = append(messages, api.Message{
				Role:    "system",
				Content: system,
			})
		}
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

			messages = append(messages, api.Message{
				Role:    "user",
				Content: scanner.Text(),
			})
		}
		skipInput = false

		output.Reset()
		agentSpoke = false
		agentFinishedThinking = false
		pendingToolCalls := []api.ToolCall{}
		err = client.Chat(ctx, &api.ChatRequest{
			Model:    model,
			Messages: messages,
			Tools:    apiTools,
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
						fmt.Print("\u001b[93mAgent\u001b[0m: ")
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
