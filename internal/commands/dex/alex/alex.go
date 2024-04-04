package alex

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/internal/common"
	"github.com/hashhavoc/teller/pkg/api/alex"
	"github.com/urfave/cli/v2"
)

func CreateAlexCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "alex",
		Usage: "get alex dex pairs",
		Action: func(c *cli.Context) error {
			resp, err := props.AlexClient.GetPairs()
			if err != nil {
				return err
			}
			dataRows := generateTableData(resp)

			headers := []string{"Base", "Target", "Last Price", "Liquidity in USD", "Contract ID"}

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
			}

			if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
				return err
			}
			return nil
		},
	}
}

func generateTableData(pairs []alex.CurrencyPair) []common.TableData {
	var dataRows []common.TableData
	for _, d := range pairs {
		row := common.TableData{
			d.Base,
			d.Target,
			strconv.FormatFloat(d.LastPrice, 'f', -1, 64),
			strconv.FormatFloat(d.LiquidityInUSD, 'f', 3, 64),
			d.BaseCurrency,
		}
		dataRows = append(dataRows, row)
	}
	return dataRows
}
