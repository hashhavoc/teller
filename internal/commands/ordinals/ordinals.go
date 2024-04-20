package ordinals

import (
	"github.com/hashhavoc/teller/internal/commands/ordinals/runes"
	"github.com/hashhavoc/teller/internal/commands/props"

	"github.com/urfave/cli/v2"
)

func CreateOrdinalsCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:    "ordinals",
		Aliases: []string{"ord"},
		Usage:   "Provides interactions with ordinals",
		Subcommands: []*cli.Command{
			runes.CreateRunesCommand(props),
		},
	}
}
