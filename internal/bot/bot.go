package bot

import (
	"fmt"
	"strings"
	"sync"

	"github.com/SomeSuperCoder/superai/internal/text"
	"github.com/SomeSuperCoder/superai/internal/util"

	"github.com/revrost/go-openrouter"
	log "github.com/sirupsen/logrus"
)

type Bot struct {
	Agents map[string]Agent
	History []openrouter.ChatCompletionMessage
}

func (b Bot) Pipeline(request string) string {
	b.Agents = loadAgents()

	// Create individual history slices
	for _, agent := range b.Agents {
		var summary = agentsSummary(b.Agents, agent.Identificator)
		var system = util.MakeSystemMessage(fmt.Sprintf(text.AGENT_BASE_PROMPT, summary, agent.Prompt))
		b.History = append(b.History, util.MakeUserMessage(request))
		agent.HistorySlice = append([]openrouter.ChatCompletionMessage{system}, b.History...)
		b.Agents[agent.Identificator] = agent
	}

	b.pipelineCycle()

	var wg sync.WaitGroup
	wg.Add(len(b.Agents))
	var results = make(chan AgentResult, len(b.Agents))

	// Summary
	for _, agent := range b.Agents {
		agent.HistorySlice = append(agent.HistorySlice, util.MakeUserMessage(text.SUMMARY_PROMPT))
		agent.Query(&wg, results)
	}

	wg.Wait()
	close(results)

	var resultsSlice []AgentResult

	for result := range results {
		resultsSlice = append(resultsSlice, result)
	}

	// Crafting the final response
	var crafter = util.MakeUserMessage(fmt.Sprintf(text.CRAFT_PROMPT, b.summarizeResults(resultsSlice)))
	var historySlice = append(b.History, crafter)
	var responseChannel = make(chan string)

	go util.Prompt("moonshotai/kimi-k2:free", responseChannel, historySlice)

	var finalResultBuilder strings.Builder

	for token := range responseChannel {
		finalResultBuilder.WriteString(token)
	}

	var finalResult = finalResultBuilder.String()

	fmt.Printf("%v\n", b.History)

	b.History = append(b.History, util.MakeAssistantMessage(finalResult))

	return finalResult
}

func (b Bot) pipelineCycle() {
	// Count live agents
	var liveCount = 0

	for _, agent := range b.Agents {
		if agent.Live {
			liveCount++
		}
	}

	if liveCount == 0 {
		log.Info("All agents sleeping - finshing the pipeline!")
		return
	}

	var wg sync.WaitGroup
	wg.Add(liveCount)

	var results = make(chan AgentResult, liveCount)

	for _, agent := range b.Agents {
		if !agent.Live { continue }

		go agent.Query(&wg, results);
	}

	wg.Wait()
	close(results)

	var mail = make(map[string][]IncomingMessage)

	for result := range results {
		for _, message := range result.Messages {
			mail[message.To] = append(mail[message.To], IncomingMessage{result.Identificator, message.Content})
		}
	}

	// TODO: handle votes

	for _, agent := range b.Agents {
		var messages = mail[agent.Identificator]
		var receiver = agent.Identificator

		// Suspend them
		if len(messages) == 0 {
			var target = b.Agents[receiver]
			target.Live = false
			b.Agents[receiver] = target
			continue
		}

		// Parse messages
		var builder strings.Builder
		builder.WriteString("# Incoming messages:\n")

		for _, message := range messages {
			builder.WriteString(fmt.Sprintf("New message from %v: ```%v```\n", message.From, message.Content))
		}

		// Send messages and move 'em to the new team
		var target = b.Agents[receiver]
		target.HistorySlice = append(b.Agents[receiver].HistorySlice, util.MakeUserMessage(builder.String()))
		target.Live = true
		b.Agents[receiver] = target
		log.Info(fmt.Sprintf("Moving %v to the next interation", target.Identificator))

	}

	b.pipelineCycle()
}

func (b Bot) summarizeResults(results []AgentResult) string {
	var builder strings.Builder

	for _, result := range results {
		builder.WriteString(fmt.Sprintf("--- BEGIN RESPONSE BY { ID: %v ; ABOUT: `%v` }", result.Identificator, b.Agents[result.Identificator].Description))
		builder.WriteString(result.Response + "\n")
		builder.WriteString("--- END ---\n")
	}

	return builder.String()
}

