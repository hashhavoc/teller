package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/hashhavoc/teller/internal/commands/bob"
	"github.com/hashhavoc/teller/internal/commands/conf"
	"github.com/hashhavoc/teller/internal/commands/contract"
	"github.com/hashhavoc/teller/internal/commands/dex"
	"github.com/hashhavoc/teller/internal/commands/names"
	"github.com/hashhavoc/teller/internal/commands/ordinals"
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/internal/commands/token"
	"github.com/hashhavoc/teller/internal/commands/transactions"
	"github.com/hashhavoc/teller/internal/commands/wallet"
	"github.com/hashhavoc/teller/internal/config"
	"github.com/hashhavoc/teller/pkg/api/alex"
	"github.com/hashhavoc/teller/pkg/api/gobob"
	"github.com/hashhavoc/teller/pkg/api/hiro"
	"github.com/hashhavoc/teller/pkg/api/ord"
	"github.com/hashhavoc/teller/pkg/api/stxtools"
	"github.com/phuslu/log"

	"github.com/urfave/cli/v2"
)

func CreateApp(glog log.Logger, version string) *cli.App {
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
	ordClient := ord.NewAPIClient(config.Endpoints.Ord)
	gobobClient := gobob.NewAPIClient(config.Endpoints.Bob)
	props := &props.AppProps{
		HeroClient:     hiroClient,
		AlexClient:     alexClient,
		StxToolsClient: stxtoolsClient,
		OrdClient:      ordClient,
		BobClient:      gobobClient,
		Config:         config,
		Logger:         glog,
	}
	app := &cli.App{
		Name:                 "teller",
		Compiled:             time.Now(),
		Version:              version,
		Usage:                "interact with the stx blockchain",
		EnableBashCompletion: true,
		Suggest:              true,
		Commands: []*cli.Command{
			conf.CreateConfigCommand(props),
			bob.CreateBobCommand(props),
			contract.CreateContractsCommand(props),
			token.CreateTokenCommand(props),
			wallet.CreateWalletCommand(props),
			dex.CreateDexCommand(props),
			transactions.CreateTransactionsCommand(props),
			ordinals.CreateOrdinalsCommand(props),
			names.CreateNameCommand(props),
		},
	}
	return app
}
