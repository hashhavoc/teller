package holders

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/internal/common"
	"github.com/hashhavoc/teller/pkg/api/hiro"
	"github.com/urfave/cli/v2"
)

func CreateFungibleTokenHoldersCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "holders",
		Usage: "Retrieves holders for a token contract",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "contract",
				Usage:    "Contract address",
				Aliases:  []string{"c"},
				Required: true,
			},
			&cli.IntFlag{
				Name:    "block",
				Usage:   "Block height to query at",
				Aliases: []string{"b"},
				Value:   0,
			},
		},
		Action: func(c *cli.Context) error {
			resp, err := props.HeroClient.GetTokenHolders(c.String("contract"), c.Int("block"))
			if err != nil {
				props.Logger.Fatal().Err(err).Msg("Error getting all tokens")
			}
			dataRows := generateTableData(resp)

			headers := []string{"Address", "Balance"}

			t := common.CreateTable(headers, dataRows)

			vpTop := viewport.New(75, 1)
			vpTop.SetContent(fmt.Sprintf("Total: %s", fmt.Sprint(len(dataRows))))

			vpBottom := viewport.New(75, 1)
			vpBottom.SetContent("Press 's' to export all addresses, 'enter' to open address in explorer, 1-9 to sort by column")

			m := tableModel{
				table:          t,
				viewportBottom: vpBottom,
				viewportTop:    vpTop,
				client:         props.HeroClient,
				logger:         props.Logger,
			}

			if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
				props.Logger.Fatal().Err(err).Msg("Failed to run program")
			}
			return nil
		},
	}
}

func generateTableData(pairs hiro.ContractHoldersResponse) []common.TableData {
	var dataRows []common.TableData
	for k, d := range pairs {
		row := common.TableData{
			k,
			d,
		}
		dataRows = append(dataRows, row)
	}
	return dataRows
}
