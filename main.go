package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	readline "github.com/nyaosorg/go-readline-ny"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	msgs := []openai.ChatCompletionMessage{}
	editor := readline.Editor{
		Prompt: func() (int, error) {
			return fmt.Print("\n> ")
		},
	}
	for {
		text, err := editor.ReadLine(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		if text == "q" || text == "quit" {
			break
		}
		msgs = append(msgs, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: text})
		resp, err := client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{Model: openai.GPT3Dot5Turbo, Messages: msgs},
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Printf("\n\n%v\n", resp.Choices[0].Message.Content)
		msgs = append(msgs, resp.Choices[0].Message)
	}

	if err := json.NewEncoder(os.Stdout).Encode(msgs); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
