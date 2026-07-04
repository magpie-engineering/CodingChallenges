package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/magpie-engineering/CodingChallenges/word_count/run"
)

func Run(args []string) int {
	var app run.AppEnv
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

func fromArgs(app *run.AppEnv, args []string) error {

	fl := flag.NewFlagSet("word_count", flag.ExitOnError)
	fl.BoolVar(&app.ByteCount, "c", false, "byte count")
	fl.BoolVar(&app.LineCount, "l", false, "line count")
	fl.BoolVar(&app.WordCount, "w", false, "word count")
	fl.BoolVar(&app.CharCount, "m", false, "char count")
	if err := fl.Parse(args); err != nil {
		return err
	}
	app.Filename = fl.Arg(0)
	if app.Filename == "" {
		return nil
	}
	if !(app.ByteCount && app.LineCount && app.WordCount && app.CharCount) {
		// no flag specified so default flags
		app.ByteCount = true
		app.LineCount = true
		app.WordCount = true
	}

	return nil
}
