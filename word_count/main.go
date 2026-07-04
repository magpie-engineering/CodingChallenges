package main

import (
	"os"

	"github.com/magpie-engineering/CodingChallenges/word_count/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
