package runes

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/internal/common"
	"github.com/hashhavoc/teller/pkg/api/ord"
	"github.com/urfave/cli/v2"
)

func CreateRunesCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "runes",
		Usage: "Provides interactions with runes",
		Action: func(c *cli.Context) error {
			resp, err := props.OrdClient.GetAllRunes()
			if err != nil {
				props.Logger.Fatal().Err(err).Msg("Error getting all tokens")
			}
			dataRows := generateTableData(resp)

			headers := []string{"Name", "Symbol", "Amount", "Mints", "Premine", "Total Supply", "Premine Percent"}

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

func generateTableData(pairs []ord.Entry) []common.TableData {
	var dataRows []common.TableData
	for _, d := range pairs {
		// premine := d.Details.Premine /d.Details.
		totalSupply := d.Details.Premine + (d.Details.Mints * d.Details.Terms.Amount)
		preminePercent := float64(d.Details.Premine) / float64(totalSupply)
		preminePercent = preminePercent * 100
		row := common.TableData{
			d.Details.SpacedRune,
			d.Details.Symbol,
			fmt.Sprintf("%d", d.Details.Terms.Amount),
			fmt.Sprintf("%d", d.Details.Mints),
			fmt.Sprintf("%d", d.Details.Premine),
			fmt.Sprintf("%d", totalSupply),
			fmt.Sprintf("%f", preminePercent),
		}
		dataRows = append(dataRows, row)
	}
	return dataRows
}
