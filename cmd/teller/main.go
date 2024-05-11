package main

import (
	"os"

	"github.com/hashhavoc/teller/internal/commands"
	"github.com/phuslu/log"
)

var Version = "development"

func main() {
	logLevelEnv := os.Getenv("TELLER_LOG_LEVEL")
	logLevel := log.ErrorLevel

	switch logLevelEnv {
	case "DEBUG":
		logLevel = log.DebugLevel
	case "INFO":
		logLevel = log.InfoLevel
	case "WARN":
		logLevel = log.WarnLevel
	case "ERROR":
		logLevel = log.ErrorLevel
	case "FATAL":
		logLevel = log.FatalLevel
	}

	glog := log.Logger{
		TimeFormat: "15:04:05",
		Level:      logLevel,
		Caller:     0,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    false,
			EndWithMessage: true,
		},
	}
	app := commands.CreateApp(glog, Version)
	err := app.Run(os.Args)
	if err != nil {
		glog.Fatal().Err(err).Msg("")
	}
}
