package main

import (
	"universalScraper/product/tgbot"
	"universalScraper/runner"
)

func main() {
	tgbot.Init()
	tgbot.StartPolling()
	runner.Init()
}
