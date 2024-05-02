package bob

import (
	token "github.com/hashhavoc/teller/internal/commands/bob/tokens"
	"github.com/hashhavoc/teller/internal/commands/props"

	"github.com/urfave/cli/v2"
)

func CreateBobCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "bob",
		Usage: "Provides interactions with bob chain",
		Subcommands: []*cli.Command{
			token.CreateBobTokenCommand(props),
		},
	}
}
