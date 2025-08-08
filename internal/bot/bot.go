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
	Agents []Agent
	History []openrouter.ChatCompletionMessage
}

func (b Bot) Pipeline(request string) string {
	b.Agents = loadAgents()

	var wg sync.WaitGroup

	for _, agent := range b.Agents {
		wg.Add(1)
		var summary = agentsSummary(b.Agents, agent.Identificator)

		var system = util.MakeSystemMessage(fmt.Sprintf(text.AGENT_BASE_PROMPT, summary, agent.Prompt))
		b.History = append(b.History, util.MakeUserMessage(request))

		work := func () {
			var responseChannel = make(chan string)
			go util.Prompt("moonshotai/kimi-k2:free", responseChannel, append([]openrouter.ChatCompletionMessage{system}, b.History...))

			var builder strings.Builder

			for part := range responseChannel {
				log.Info(fmt.Sprintf("Data from <%v>", agent.Identificator))
				builder.WriteString(part)
			}
			log.Info(fmt.Sprintf("Final response from %v: %v", agent.Identificator, builder.String()))

			wg.Done()
		}
		go work()
	}

	wg.Wait()
	return ""
}

