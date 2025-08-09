package util

import (
	"context"
	"errors"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
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
		log.Fatal(err.Error())
	}

	for {
		response, err := stream.Recv()

		if err != nil && err != io.EOF {
			log.Fatal(err.Error())
		}

		if errors.Is(err, io.EOF) {
			close(channel)
			stream.Close()
			return
		}

		channel <- response.Choices[0].Delta.Content
	}
}

