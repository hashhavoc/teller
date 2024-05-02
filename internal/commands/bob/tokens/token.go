package token

import (
	"github.com/hashhavoc/teller/internal/commands/bob/tokens/ft"
	"github.com/hashhavoc/teller/internal/commands/props"

	"github.com/urfave/cli/v2"
)

func CreateBobTokenCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "token",
		Usage: "Provides interactions with bob tokens",
		Subcommands: []*cli.Command{
			ft.CreateFungibleTokensCommand(props),
		},
	}
}
