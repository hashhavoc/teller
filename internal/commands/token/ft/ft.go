package ft

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/internal/common"
	"github.com/hashhavoc/teller/pkg/api/hiro"
	"github.com/urfave/cli/v2"
)

func CreateFungibleTokensCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:    "fungible",
		Aliases: []string{"ft"},
		Usage:   "Provides interactions with fungible tokens",
		Action: func(c *cli.Context) error {
			resp, err := props.HeroClient.GetAllTokens()
			if err != nil {
				props.Logger.Fatal().Err(err).Msg("Error getting all tokens")
			}
			dataRows := generateTableData(resp)

			headers := []string{"Name", "Symbol", "Decimals", "Total Supply", "Contract ID"}

			t := common.CreateTable(headers, dataRows)

			vpTop := viewport.New(75, 1)
			vpTop.SetContent(fmt.Sprintf("Total: %s", fmt.Sprint(len(dataRows))))

			vpBottom := viewport.New(75, 1)
			vpBottom.SetContent("Press 'c' to copy address, 'a' to export all addresses, 'h' to view holders")

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

func generateTableData(pairs []hiro.TokenResult) []common.TableData {
	var dataRows []common.TableData
	for _, d := range pairs {
		row := common.TableData{
			d.Name,
			d.Symbol,
			fmt.Sprintf("%d", d.Decimals),
			d.TotalSupply,
			d.ContractPrincipal,
		}
		dataRows = append(dataRows, row)
	}
	return dataRows
}
