package dex

import (
	alexcmd "github.com/hashhavoc/teller/internal/commands/dex/alex"
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/urfave/cli/v2"
)

func CreateDexCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "dex",
		Usage: "Provides interactions with multiple dex",
		Subcommands: []*cli.Command{
			alexcmd.CreateAlexCommand(props),
		},
	}
}
