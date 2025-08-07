package main

import (
	"fmt"

	"github.com/SomeSuperCoder/superai/internal/util"
	"github.com/revrost/go-openrouter"
)

func main() {
	defer fmt.Println()

	var responseChannel = make(chan string)
	go util.Prompt("moonshotai/kimi-k2:free", responseChannel, []openrouter.ChatCompletionMessage{
		{
			Role: openrouter.ChatMessageRoleUser,
			Content: openrouter.Content{ Text: "Hello, Sup?" },
		},
	})

	for part := range responseChannel {
		fmt.Print(part)
	}
}
