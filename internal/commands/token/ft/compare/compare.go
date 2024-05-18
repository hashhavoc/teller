package compare

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/internal/common"
	"github.com/hashhavoc/teller/pkg/api/hiro"
	"github.com/hashhavoc/teller/pkg/utils/uint128"
	"github.com/urfave/cli/v2"
)

func CreateFungibleTokenHoldersCompareCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "compare",
		Usage: "Retrieves holders for a token contract and compare based on height",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "contract",
				Usage:    "Contract address",
				Aliases:  []string{"c"},
				Required: true,
			},
			&cli.IntFlag{
				Name:    "first",
				Usage:   "Block height to query at",
				Aliases: []string{"f"},
				Value:   0,
			},
			&cli.IntFlag{
				Name:    "second",
				Usage:   "Block height to query at",
				Aliases: []string{"s"},
				Value:   0,
			},
		},
		Action: func(c *cli.Context) error {
			firstResp, err := props.HeroClient.GetTokenHolders(c.String("contract"), c.Int("first"))
			if err != nil {
				props.Logger.Fatal().Err(err).Msg("Error getting all tokens")
			}

			secondResp, err := props.HeroClient.GetTokenHolders(c.String("contract"), c.Int("second"))
			if err != nil {
				props.Logger.Fatal().Err(err).Msg("Error getting all tokens")
			}

			dataRows := generateTableData(firstResp, secondResp)

			headers := []string{"Address", "First", "Second"}

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

type AddressData struct {
	Before     uint128.Uint128
	After      uint128.Uint128
	Difference uint128.Uint128
}

func generateTableData(firstResp hiro.ContractHoldersResponse, secondResp hiro.ContractHoldersResponse) []common.TableData {
	var dataRows []common.TableData

	// Create a map to store addresses and their data
	addressData := make(map[string]AddressData)

	// Process the first response
	for address, amount := range firstResp {
		value, err := uint128.FromString(amount)
		if err != nil {
			continue
		}
		addressData[address] = AddressData{Before: value}
	}

	// Process the second response and calculate differences
	for address, amount := range secondResp {
		value, err := uint128.FromString(amount)
		if err != nil {
			continue
		}
		if data, exists := addressData[address]; exists {
			var difference uint128.Uint128
			direct := data.Before.Cmp(value)
			if direct == 1 {
				difference = data.Before.Sub(value)
			} else if direct == -1 {
				difference = value.Sub(data.Before)
			} else {
				difference = uint128.Zero
			}
			addressData[address] = AddressData{
				Before:     data.Before,
				After:      value,
				Difference: difference,
			}
		} else {
			addressData[address] = AddressData{
				Before:     uint128.From64(0),
				After:      value,
				Difference: value,
			}
		}
	}

	// Convert the map to a slice of TableData
	for address, data := range addressData {
		row := common.TableData{
			address,
			data.Before.String(),
			data.After.String(),
		}

		dataRows = append(dataRows, row)
	}

	return dataRows
}
