package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/SomeSuperCoder/superai/internal/bot"
)

func main() {
	var master bot.Bot

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(" >> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var result = master.Pipeline(input)
		fmt.Println("==============================================")
		fmt.Println(result)
	}
}

