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

type AgentResult struct {
	Identificator string
	Response string `yaml:"response"`
	Messages []OutgoingMessage `yaml:"messages"`
	Continue bool `yaml:"continue"`
}

type OutgoingMessage struct {
	To string
	Content string
}

type IncomingMessage struct {
	From string
	Content string
}

type Agent struct {
	// Based on filenames
	Identificator string
	// Parsed data
	Description string `yaml:"description"`
	Prompt string `yaml:"prompt"`
	// Defined at runtime
	HistorySlice []openrouter.ChatCompletionMessage
	Live bool
}

func (a Agent) Query(wg *sync.WaitGroup, results chan AgentResult) string {
	defer wg.Done()

	var responseChannel = make(chan string)
	go util.Prompt("moonshotai/kimi-k2:free", responseChannel, a.HistorySlice)

	var builder strings.Builder

	for part := range responseChannel {
		fmt.Println(a.Identificator)
		builder.WriteString(part)
	}

	var final = builder.String()

	// Parsing
	var result AgentResult
	yaml.Unmarshal([]byte(final), &result)
	result.Identificator = a.Identificator

	results <- result

	return final
}

func loadAgents() map[string]Agent {
	entries, err := os.ReadDir("agents")

	if err != nil {
		log.Fatal(err.Error())
	}

	var result = make(map[string]Agent)

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
		agent.Live = true

		result[agent.Identificator] = agent
	}

	return result
}

func agentsSummary(agents map[string]Agent, exculde string) string {
	var builder strings.Builder

	var i = 0

	for _, agent := range agents {
		if agent.Identificator != exculde {
			builder.WriteString(fmt.Sprintf("%v) ID: %v ; Description: %v\n", i + 1, agent.Identificator, agent.Description))
			i++
		}
	}

	return builder.String()
}

