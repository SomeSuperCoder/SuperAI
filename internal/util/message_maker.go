package util

import "github.com/revrost/go-openrouter"

func makeMessage(text string, role string) openrouter.ChatCompletionMessage {
	return openrouter.ChatCompletionMessage{
		Role: role,
		Content: openrouter.Content{ Text: text },
	}
}

func MakeUserMessage(text string) openrouter.ChatCompletionMessage {
	return makeMessage(text, openrouter.ChatMessageRoleUser)
}

func MakeSystemMessage(text string) openrouter.ChatCompletionMessage {
	return makeMessage(text, openrouter.ChatMessageRoleSystem)
}

func MakeAssistantMessage(text string) openrouter.ChatCompletionMessage {
	return makeMessage(text, openrouter.ChatMessageRoleAssistant)
}

