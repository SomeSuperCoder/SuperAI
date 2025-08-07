package util

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	openrouter "github.com/revrost/go-openrouter"
)

var client = openrouter.NewClient(
	os.Getenv("SUPER_AI_API_KEY"),
)

func Prompt(model string, channel chan string, messages []openrouter.ChatCompletionMessage) {
	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
			Stream:   true,
		},
	)

	if err != nil {
		log.Panicln(err.Error())
	}

	for {
		response, err := stream.Recv()

		if err != nil && err != io.EOF {
			log.Panicln(err.Error())
		}

		if errors.Is(err, io.EOF) {
			stream.Close()
			close(channel)
			return
		}

		channel <- response.Choices[0].Delta.Content
	}
}

