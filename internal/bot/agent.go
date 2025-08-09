package bot

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/SomeSuperCoder/superai/internal/util"

	"github.com/revrost/go-openrouter"
	"gopkg.in/yaml.v2"
	log "github.com/sirupsen/logrus"
)

type AgentResoponse struct {
	Id string
	Text string
}

type AgentResult struct {
	Response string `yaml:"response"`
	Messages []AgentMessage `yaml:"messages"`
	Continue bool `yaml:"continue"`
}

type AgentMessage struct {
	To string
	Content string
}

type Agent struct {
	Identificator string
	Description string `yaml:"description"`
	Prompt string `yaml:"prompt"`
}

func (a Agent) Query(system openrouter.ChatCompletionMessage, messages []openrouter.ChatCompletionMessage, wg *sync.WaitGroup, responses chan AgentResoponse) string {
	defer wg.Done()

	var responseChannel = make(chan string)
	go util.Prompt("moonshotai/kimi-k2:free", responseChannel, append([]openrouter.ChatCompletionMessage{system}, messages...))

	var builder strings.Builder

	for part := range responseChannel {
		log.Info(fmt.Sprintf("Data from <%v>", a.Identificator))
		builder.WriteString(part)
	}

	var final = builder.String()

	responses <- AgentResoponse{a.Identificator, final}

	log.Info(fmt.Sprintf("Final response from %v", a.Identificator))
	return final
}

func loadAgents() []Agent {
	entries, err := os.ReadDir("agents")

	if err != nil {
		log.Fatal(err.Error())
	}

 	var result = make([]Agent, 0, len(entries))

	for _, entry := range entries {
		log.Info("Loading agent: " + entry.Name() + "...")

		data, err := os.ReadFile("agents/" + entry.Name())

		if err != nil {
			log.Fatal(err.Error())
		}

		var agent Agent
		err = yaml.Unmarshal(data, &agent)

		if err != nil {
			log.Fatal(err.Error())
		}

		agent.Identificator = entry.Name()

		result = append(result, agent)
	}

	return result
}

func agentsSummary(agents []Agent, exculde string) string {
	var builder strings.Builder

	for i, agent := range agents {
		if agent.Identificator != exculde {
			builder.WriteString(fmt.Sprintf("%v) ID: %v ; Description: %v\n", i+1, agent.Identificator, agent.Description))
		}
	}

	return builder.String()
}

