package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/SomeSuperCoder/superai/internal/util"

	openrouter "github.com/revrost/go-openrouter"
)

func main() {
	client := openrouter.NewClient(
		os.Getenv("SUPER_AI_API_KEY"),
	)

	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Model: "meta-llama/llama-3.3-70b-instruct:free",
			Messages: []openrouter.ChatCompletionMessage{
				{
					Role: openrouter.ChatMessageRoleUser,
					Content: openrouter.Content{Text: "Hello! Sup?"},
				},
			},
			Stream: true,
		},
	)

	if err != nil {
		log.Panicln(err.Error())
	}

	defer stream.Close()
	defer fmt.Println()

	for {
		response, err := stream.Recv()

		if err != nil && err != io.EOF {
			fmt.Println(err.Error())
		}
		if errors.Is(err, io.EOF) {
			return
		}

		data, err := json.MarshalIndent(response, "", " ")

		if err != nil {
			log.Panicln(err.Error())
		}

		var parsed_response util.Response

		err = json.Unmarshal(data, &parsed_response)

		if err != nil {
			log.Panicln(err.Error())
		}

		fmt.Print(parsed_response.Choices[0].Delta.Content)
	}
}

