package main

import (
	"os"

	"github.com/hashhavoc/teller/internal/commands"
	"github.com/phuslu/log"
)

func main() {
	glog := log.Logger{
		TimeFormat: "15:04:05",
		Caller:     0,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    false,
			EndWithMessage: true,
		},
	}
	app := commands.CreateApp(glog)
	err := app.Run(os.Args)
	if err != nil {
		glog.Fatal().Err(err).Msg("")
	}
}
