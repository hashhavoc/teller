package token

import (
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/internal/commands/token/ft"
	"github.com/hashhavoc/teller/internal/commands/token/ft/holders"
	"github.com/hashhavoc/teller/internal/commands/token/nft"

	"github.com/urfave/cli/v2"
)

func CreateTokenCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "token",
		Usage: "Provides interactions with tokens",
		Subcommands: []*cli.Command{
			nft.CreateNonFungibleTokensCommand(props),
			ft.CreateFungibleTokensCommand(props),
			holders.CreateFungibleTokenHoldersCommand(props),
		},
	}
}
