package main

import (
	"os"

	"github.com/magpie-engineering/CodingChallenges/redisserver/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
