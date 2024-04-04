package nft

import (
	"fmt"
	"strings"

	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/urfave/cli/v2"
)

func CreateNonFungibleTokensCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:    "nonfungible",
		Aliases: []string{"nft"},
		Usage:   "Provides interactions with non-fungible tokens",
		Subcommands: []*cli.Command{
			createHoldingsCommand(props),
		},
	}
}

func createHoldingsCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "holdings",
		Usage: "Retrieves holdings for a specific non-fungible token",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "principal",
				Usage:    "Principal address",
				Aliases:  []string{"p"},
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			resp, err := props.HeroClient.GetNFTHoldings(c.String("principal"))
			if err != nil {
				props.Logger.Fatal().Err(err).Msg("Error getting nft holdings")
			}
			for _, nft := range resp {
				split := strings.Split(nft.AssetIdentifier, "::")
				fmt.Println(split[1])
			}
			return nil
		},
	}
}
