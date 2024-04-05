package initialize

import (
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/pkg/api/alex"
	"github.com/hashhavoc/teller/pkg/api/hiro"
	"github.com/hashhavoc/teller/pkg/api/stxtools"
	"github.com/urfave/cli/v2"
)

// CreateContractsCommand creates the contracts command and its subcommands.
func CreateInitCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Creates a new configuration file",
		Action: func(c *cli.Context) error {
			props.Config.Endpoints.Hiro = hiro.DefaultApiBase
			props.Config.Endpoints.Alex = alex.DefaultApiBase
			props.Config.Endpoints.StxTools = stxtools.DefaultApiBase
			props.Config.WriteConfig()
			return nil
		},
	}
}
