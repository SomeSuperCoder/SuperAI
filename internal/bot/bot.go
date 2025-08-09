package bot

import (
	"fmt"
	"sync"

	"github.com/SomeSuperCoder/superai/internal/text"
	"github.com/SomeSuperCoder/superai/internal/util"

	"github.com/revrost/go-openrouter"
	log "github.com/sirupsen/logrus"
)

type Bot struct {
	Agents []Agent
	History []openrouter.ChatCompletionMessage
}

func (b Bot) Pipeline(request string) string {
	b.Agents = loadAgents()
	agentsAmount := len(b.Agents)

	var wg sync.WaitGroup
	wg.Add(agentsAmount)

	var responses = make(chan AgentResoponse, agentsAmount)

	for _, agent := range b.Agents {
		var summary = agentsSummary(b.Agents, agent.Identificator)

		var system = util.MakeSystemMessage(fmt.Sprintf(text.AGENT_BASE_PROMPT, summary, agent.Prompt))
		b.History = append(b.History, util.MakeUserMessage(request))

		go agent.Query(system, b.History, &wg, responses);
	}

	wg.Wait()
	close(responses)

	for response := range responses {
		log.Infoln(fmt.Sprintf("Final response from %v: %v", response.Id, response.Text))
	}

	return ""
}

