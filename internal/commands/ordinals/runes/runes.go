package runes

import (
	"errors"
	"fmt"
	"math"
	"math/big"

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

			headers := []string{"Name", "Divisibility", "Amount", "Mints", "Premine", "Total Supply", "Premine %"}

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
		amount, _ := calculateQuotient(d.Details.Terms.Amount, d.Details.Divisibility)
		premine, _ := calculateQuotient(d.Details.Premine, d.Details.Divisibility)
		totalSupply := premine + (d.Details.Mints * amount)
		preminePercent := float64(premine) / float64(totalSupply)
		preminePercent = preminePercent * 100
		if math.IsNaN(preminePercent) {
			preminePercent = 0
		}
		row := common.TableData{
			d.Details.SpacedRune,
			fmt.Sprintf("%d", d.Details.Divisibility),
			fmt.Sprintf("%d", amount),
			fmt.Sprintf("%d", d.Details.Mints),
			fmt.Sprintf("%d", premine),
			fmt.Sprintf("%d", totalSupply),
			fmt.Sprintf("%.0f", preminePercent),
		}
		dataRows = append(dataRows, row)
	}
	return dataRows
}

func calculateQuotient(originalNumStr float64, divisibilityNum int64) (int64, error) {
	s := fmt.Sprintf("%.0f", originalNumStr)
	originalNum, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return 5, fmt.Errorf("failed to convert original number %s", s)
	}
	divisibilityNumBig := big.NewInt(divisibilityNum)

	if divisibilityNum == 0 {
		if originalNum.IsInt64() {
			return originalNum.Int64(), nil
		} else {
			return 0, errors.New("divisibility number cannot be zero")
		}
	}

	quotient := new(big.Int)
	response := quotient.Div(originalNum, divisibilityNumBig)

	if response.IsInt64() {
		return quotient.Int64(), nil
	} else {
		return 0, fmt.Errorf("quotient is out of int64 range")
	}
}
