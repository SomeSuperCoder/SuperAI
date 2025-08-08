package main

import (
	"fmt"

	"github.com/SomeSuperCoder/superai/internal/bot"
)

func main() {
	defer fmt.Println()

	var master bot.Bot

	master.Pipeline("Is ArchLinux worth it?")
}
