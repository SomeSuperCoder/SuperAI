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

	for _, agent := range b.Agents {
		// Create an idividual history slice
		var summary = agentsSummary(b.Agents, agent.Identificator)
		var system = util.MakeSystemMessage(fmt.Sprintf(text.AGENT_BASE_PROMPT, summary, agent.Prompt))
		b.History = append(b.History, util.MakeUserMessage(request))
		agent.HistorySlice = append([]openrouter.ChatCompletionMessage{system}, b.History...)
		b.Agents[agent.Identificator] = agent
	}

	b.pipelineCycle(b.Agents)

	return ""
}

func (b Bot) pipelineCycle(liveAgents map[string]Agent) {
	if len(liveAgents) == 0 {
		log.Info("All agents sleeping - finshing the pipeline!")
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(liveAgents))

	var results = make(chan AgentResult, len(liveAgents))

	for _, agent := range liveAgents {
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

	var newTeam = make(map[string]Agent)

	for receiver, messages := range mail {
		// Parse messages
		var builder strings.Builder
		builder.WriteString("# Incmoing messages:\n")

		for _, message := range messages {
			builder.WriteString(fmt.Sprintf("New message from %v: ```%v```\n", message.From, message.Content))
		}

		// Send messages and move 'em to the new team
		var target = liveAgents[receiver]
		target.HistorySlice = append(liveAgents[receiver].HistorySlice, util.MakeUserMessage(builder.String()))
		target.Identificator = receiver // Fix later!
		newTeam[receiver] = target
		log.Info(fmt.Sprintf("Moving %v to the next interation", target.Identificator))
	}

	b.pipelineCycle(newTeam)
}

