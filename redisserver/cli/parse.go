package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/magpie-engineering/CodingChallenges/redisserver/redis"
)

func Run(args []string) int {
	var app redis.AppEnv
	err := fromArgs(&app, args)
	if err != nil {
		return 2
	}

	if err = app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error:%v\n", err)
		return 1
	}

	return 0
}

func fromArgs(app *redis.AppEnv, args []string) error {

	fl := flag.NewFlagSet("redis", flag.ExitOnError)
	fl.IntVar(&app.Port, "p", 6379, "port")
	fl.StringVar(&app.Addr, "addr", "127.0.0.1", "bind address")
	if err := fl.Parse(args); err != nil {
		return err
	}
	return nil
}
