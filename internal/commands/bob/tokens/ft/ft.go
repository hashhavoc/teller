package ft

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/internal/common"
	"github.com/hashhavoc/teller/pkg/api/gobob"
	"github.com/urfave/cli/v2"
)

func CreateFungibleTokensCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:    "fungible",
		Aliases: []string{"ft"},
		Usage:   "Provides interactions with fungible tokens",
		Action: func(c *cli.Context) error {
			resp, err := props.BobClient.GetAllTokens()
			if err != nil {
				props.Logger.Fatal().Err(err).Msg("Error getting all tokens")
			}
			dataRows := generateTableData(resp)

			headers := []string{"Name", "Symbol", "Decimals", "Total Supply", "Contract ID", "Holders", "Type"}

			t := common.CreateTable(headers, dataRows)

			vpTop := viewport.New(75, 1)
			vpTop.SetContent(fmt.Sprintf("Total: %s", fmt.Sprint(len(dataRows))))

			vpBottom := viewport.New(75, 1)
			vpBottom.SetContent("Press 'a' to export all addresses, 'h' to view holders")

			m := tableModel{
				table:          t,
				viewportBottom: vpBottom,
				viewportTop:    vpTop,
				client:         props.BobClient,
				logger:         props.Logger,
			}

			if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
				props.Logger.Fatal().Err(err).Msg("Failed to run program")
			}
			return nil
		},
	}
}

func generateTableData(pairs []gobob.TokenItems) []common.TableData {
	var dataRows []common.TableData
	for _, d := range pairs {
		row := common.TableData{
			d.Name,
			d.Symbol,
			d.Decimals,
			d.TotalSupply,
			d.Address,
			d.Holders,
			d.Type,
		}
		dataRows = append(dataRows, row)
	}
	return dataRows
}
