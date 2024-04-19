package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/hashhavoc/teller/internal/commands/contract"
	"github.com/hashhavoc/teller/internal/commands/dex"
	"github.com/hashhavoc/teller/internal/commands/initialize"
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/internal/commands/token"
	"github.com/whoabuddy/teller/internal/commands/transactions"
	"github.com/hashhavoc/teller/internal/commands/wallet"
	"github.com/hashhavoc/teller/internal/config"
	"github.com/hashhavoc/teller/pkg/api/alex"
	"github.com/hashhavoc/teller/pkg/api/hiro"
	"github.com/hashhavoc/teller/pkg/api/stxtools"
	"github.com/phuslu/log"

	"github.com/urfave/cli/v2"
)

func CreateApp(glog log.Logger) *cli.App {
	dirname, err := os.UserHomeDir()
	if err != nil {
		glog.Fatal().Err(err).Msg("Failed to get user home directory")
	}
	configPath := fmt.Sprintf("%s/.teller.yaml", dirname)
	config := config.NewConfig(configPath)

	err = config.ReadConfig()
	if err != nil {
		glog.Debug().Err(err).Msg("Failed to read config")
	}

	hiroClient := hiro.NewAPIClient(config.Endpoints.Hiro)
	alexClient := alex.NewAPIClient(config.Endpoints.Alex)
	stxtoolsClient := stxtools.NewAPIClient(config.Endpoints.StxTools)
	props := &props.AppProps{
		HeroClient:     hiroClient,
		AlexClient:     alexClient,
		StxToolsClient: stxtoolsClient,
		Config:         config,
		Logger:         glog,
	}
	app := &cli.App{
		Name:                 "teller",
		Compiled:             time.Now(),
		Version:              "v0.0.1",
		Usage:                "interact with the stx blockchain",
		EnableBashCompletion: true,
		Suggest:              true,
		Commands: []*cli.Command{
			initialize.CreateInitCommand(props),
			contract.CreateContractsCommand(props),
			token.CreateTokenCommand(props),
			wallet.CreateWalletCommand(props),
			dex.CreateDexCommand(props),
			transactions.CreateTransactionsCommand(props),
		},
	}
	return app
}
