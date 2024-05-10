package conf

import (
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/urfave/cli/v2"
)

func CreateConfigCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Commands to manage the configuration file",
		Subcommands: []*cli.Command{
			CreateInitCommand(props),
			CreateSetCommand(props),
		},
	}
}

func CreateInitCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Creates a new configuration file",
		Action: func(c *cli.Context) error {
			props.Config.WriteConfig()
			return nil
		},
	}
}

func CreateSetCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "Creates a new configuration file",
		Subcommands: []*cli.Command{
			CreateSetEndpointCommand(props),
		},
	}
}

func CreateSetEndpointCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "endpoint",
		Usage: "set endpoint for the config",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "hiro",
				Usage:    "Hiro API Base URL",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "ord",
				Usage:    "Ordinals ORD API Base URL",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "alex",
				Usage:    "ALEX DEX API Base URL",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "stxtools",
				Usage:    "STXTools API Base URL",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "bob",
				Usage:    "GOBOB API Base URL",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {

			if c.String("hiro") != "" {
				props.Config.Endpoints.Hiro = c.String("hiro")
			}
			if c.String("ord") != "" {
				props.Config.Endpoints.Ord = c.String("ord")
			}
			if c.String("alex") != "" {
				props.Config.Endpoints.Alex = c.String("alex")
			}
			if c.String("stxtools") != "" {
				props.Config.Endpoints.StxTools = c.String("stxtools")
			}
			if c.String("bob") != "" {
				props.Config.Endpoints.Bob = c.String("bob")
			}
			err := props.Config.WriteConfig()
			if err != nil {
				props.Logger.Fatal().Err(err).Msg("Error writing config")
			}
			return nil
		},
	}
}
