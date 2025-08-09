package main

import (
	"fmt"

	"github.com/SomeSuperCoder/superai/internal/bot"
)

func main() {
	defer fmt.Println()

	var master bot.Bot

	master.Pipeline("How many r's are there in 'ausdgruaysdyasgdrasdygugrrasdygr123'?")
}

