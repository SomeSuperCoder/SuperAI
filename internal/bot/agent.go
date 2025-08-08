package bot

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/revrost/go-openrouter"
	log "github.com/sirupsen/logrus"
)

type Agent struct {
	Identificator string
	Description string `yaml:"description"`
	Prompt string `yaml:"prompt"`
}

func (a Agent) Query(messages []openrouter.ChatCompletionMessage) {

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

