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

			headers := []string{"Name", "Divisibility", "Amount", "Mints", "Premine", "Total Supply", "Cap", "Premine %", "Block", "Terms"}

			t := common.CreateTable(headers, dataRows)

			vpTop := viewport.New(75, 1)
			vpTop.SetContent(fmt.Sprintf("Total: %s", fmt.Sprint(len(dataRows))))

			vpBottom := viewport.New(75, 1)
			vpBottom.SetContent("Press 'c' to copy name, 'a' to export all runes")

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
		var preminePercent float64

		remaining := (d.Details.Terms.Cap - d.Details.Mints) * d.Details.Terms.Amount
		totalSupply := (d.Details.Premine + remaining + (d.Details.Mints * d.Details.Terms.Amount))

		if totalSupply != 0 {
			fmt.Print(float64(d.Details.Premine) / float64(totalSupply))
			preminePercent = float64(d.Details.Premine) / float64(totalSupply) * 100
		} else {
			preminePercent = 0 * 100
		}
		row := common.TableData{
			d.Details.SpacedRune,
			fmt.Sprintf("%d", d.Details.Divisibility),
			fmt.Sprintf("%d", d.Details.Terms.Amount),
			fmt.Sprintf("%d", d.Details.Mints),
			fmt.Sprintf("%d", d.Details.Premine),
			fmt.Sprintf("%d", totalSupply),
			fmt.Sprintf("%d", d.Details.Terms.Cap),
			fmt.Sprintf("%.0f", preminePercent),
			fmt.Sprintf("%d", d.Details.Block),
			fmt.Sprintf("%v", d.Details.TermsEnabled),
		}
		dataRows = append(dataRows, row)
	}
	return dataRows
}
