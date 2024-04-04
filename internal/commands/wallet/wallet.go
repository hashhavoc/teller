package wallet

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/internal/common"
	"github.com/hashhavoc/teller/pkg/api/hiro"
	"github.com/hashhavoc/teller/pkg/utils"
	"github.com/hashhavoccat/stacks-go"
	"github.com/urfave/cli/v2"
)

func CreateWalletCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "wallet",
		Usage: "Provides interactions with wallets",
		Subcommands: []*cli.Command{
			createBalanceCommand(props),
			createBalancesCommand(props),
			createAddWalletCommand(props),
			createRemoveWalletCommand(props),
			createGenerateWalletCommand(props),
		},
	}
}

// CreateGreetCommand creates a new CLI command for greeting
func createGenerateWalletCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "gen",
		Usage: "Create a endless amounts of STX wallet addresses",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "networkType",
				Aliases: []string{"n"},
				Usage:   "Specify the network type",
				Value:   "mainnet",
			},
			&cli.IntFlag{
				Name:     "amount",
				Aliases:  []string{"a"},
				Usage:    "Specify the amount",
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "private",
				Aliases: []string{"p"},
				Usage:   "Specify if you want to return a hex version of the private key",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "Specify the csv file location to output contents to (default: stdout)",
				Value:   "",
			},
		},
		Action: func(c *cli.Context) error {
			networkType := c.String("networkType")
			outputFile := c.String("file")
			amount := c.Int("amount")
			private := c.Bool("private")

			var networkVersion int
			switch networkType {
			case "mainnet":
				networkVersion = stacks.MainnetSingleSig
			case "testnet":
				networkVersion = stacks.TestnetSingleSig
			default:
				fmt.Println("Invalid network type. Please specify either 'mainnet' or 'testnet'.")
				return nil
			}

			var addresses []Address
			for i := 0; i < amount; i++ {
				privKey, err := utils.MakeRandomPrivKey()
				if err != nil {
					props.Logger.Err(err).Msg("Error generating private key")
					continue
				}

				encoded := hex.EncodeToString(privKey.PrivateKey.Serialize())

				addr, err := stacks.NewAddress(encoded, networkVersion)
				if err != nil {
					props.Logger.Err(err).Msg("Error generating address")
					continue
				}
				addresses = append(addresses, Address{PrivateKey: encoded, Address: addr})
			}

			// Check if outputFile is specified and open file if necessary
			var file *os.File
			var err error
			if outputFile != "" {
				file, err = os.Create(outputFile)
				if err != nil {
					props.Logger.Err(err).Msg("Failed to create file")
				}
				defer file.Close()
			}

			for _, addr := range addresses {
				if private {
					if file != nil {
						// Write to file
						_, err := file.WriteString(fmt.Sprintf("%s,%s\n", addr.Address, addr.PrivateKey))
						if err != nil {
							return fmt.Errorf("failed to write to file: %v", err)
						}
					} else {
						// Print to stdout
						fmt.Printf("%s,%s\n", addr.Address, addr.PrivateKey)
					}
				} else {
					if file != nil {
						// Write to file
						_, err := file.WriteString(fmt.Sprintf("%s\n", addr.Address))
						if err != nil {
							return fmt.Errorf("failed to write to file: %v", err)
						}
					} else {
						// Print to stdout
						fmt.Println(addr.Address)
					}
				}
			}

			return nil
		},
	}
}

func createAddWalletCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "add wallets to the config",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "principal",
				Usage:    "Principal address",
				Aliases:  []string{"p"},
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			err := props.Config.AddWallet(c.String("principal"))
			if err != nil {
				props.Logger.Fatal().Err(err).Msg("Error adding wallet")
			}
			err = props.Config.WriteConfig()
			if err != nil {
				props.Logger.Fatal().Err(err).Msg("Error writing config")
			}
			return nil
		},
	}
}

func createRemoveWalletCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "remove wallets to the config",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "principal",
				Usage:    "Principal address",
				Aliases:  []string{"p"},
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			props.Config.RemoveWallet(c.String("principal"))
			err := props.Config.WriteConfig()
			if err != nil {
				props.Logger.Fatal().Err(err).Msg("Error writing config")
			}
			return nil
		},
	}
}

func createBalanceCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "balance",
		Usage: "view balance for a specific address",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "principal",
				Usage:    "Principal address",
				Aliases:  []string{"p"},
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			var stxBalance int64
			var rows []table.Row

			fungibleTokenBalances := make(map[string]int64)
			nonFungibleTokenCounts := make(map[string]int64)

			address := c.String("principal")
			fmt.Println("Fetching balance for address:", address)
			resp, err := props.HeroClient.GetAccountBalance(address)
			if err != nil {
				return err
			}

			// Convert the STX balance for simplicity
			stxBalance, err = strconv.ParseInt(resp.Stx.Balance, 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing STX balance: %v", err)
			}

			// Prepare the table
			headers := []table.Column{
				{Title: "Name", Width: len("Name")},
				{Title: "Type", Width: len("Type")},
				{Title: "Balance", Width: len("Balance")},
				{Title: "Contract ID", Width: len("Contract ID")},
				{Title: "Display Name", Width: len("Display Name")},
			}

			maxWidths := make([]int, len(headers))
			for i, header := range headers {
				maxWidths[i] = header.Width
			}

			rows = append(rows, table.Row{"stx", "STX", fmt.Sprint(stxBalance), "", ""})
			for k, v := range resp.FungibleTokens {
				balance, _ := strconv.ParseInt(v.Balance, 10, 64)
				if existingBalance, exists := fungibleTokenBalances[k]; exists {
					fungibleTokenBalances[k] = existingBalance + balance
				} else {
					fungibleTokenBalances[k] = balance
				}
			}
			for k, v := range resp.NonFungibleTokens {
				count, _ := strconv.ParseInt(v.Count, 10, 64)
				if existingCount, exists := nonFungibleTokenCounts[k]; exists {
					nonFungibleTokenCounts[k] = existingCount + count
				} else {
					nonFungibleTokenCounts[k] = count
				}
			}

			for k, balance := range fungibleTokenBalances {
				split := strings.Split(k, "::")
				contractName, err := props.HeroClient.GetContractReadOnly(split[0], "get-name", "string", []string{})
				if err != nil {
					rows = append(rows, table.Row{split[1], "Fungible", strconv.FormatInt(balance, 10), split[0], ""})
				}
				rows = append(rows, table.Row{split[1], "Fungible", strconv.FormatInt(balance, 10), split[0], strings.TrimSpace(contractName)})
			}

			for k, count := range nonFungibleTokenCounts {
				split := strings.Split(k, "::")
				rows = append(rows, table.Row{split[1], "Non-Fungible", strconv.FormatInt(count, 10), split[0], ""})
			}

			for _, row := range rows {
				for i, cell := range row {
					cellStr := fmt.Sprint(cell)
					if len(cellStr) > maxWidths[i] {
						maxWidths[i] = len(cellStr)
					}
				}
			}

			for i, maxWidth := range maxWidths {
				headers[i].Width = maxWidth
			}

			t := table.New(
				table.WithColumns(headers),
				table.WithRows(rows),
				table.WithFocused(true),
				table.WithStyles(common.TableStyles),
			)

			// Render the table
			m := tableModel{table: t, client: props.HeroClient}
			if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
				props.Logger.Fatal().Err(err).Msg("Failed to run program")
			}

			return nil
		},
	}
}

func createBalancesCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "balances",
		Usage: "view balances",
		Action: func(c *cli.Context) error {
			var wallets []hiro.BalanceResponse
			var stxBalance int64
			var rows []table.Row

			fungibleTokenBalances := make(map[string]int64)
			nonFungibleTokenCounts := make(map[string]int64)

			for _, w := range props.Config.Wallets {
				fmt.Println("Wallet:", w)
				resp, err := props.HeroClient.GetAccountBalance(w)
				if err != nil {
					return err
				}
				wallets = append(wallets, resp)
			}

			for _, w := range wallets {
				x, _ := strconv.ParseInt(w.Stx.Balance, 10, 64)
				stxBalance += x
			}

			headers := []table.Column{
				{Title: "Name", Width: len("Name")},
				{Title: "Type", Width: len("Type")},
				{Title: "Balance", Width: len("Balance")},
				{Title: "Contract ID", Width: len("Contract ID")},
				{Title: "Display Name", Width: len("Display Name")},
			}

			maxWidths := make([]int, len(headers))
			for i, header := range headers {
				maxWidths[i] = header.Width
			}

			rows = append(rows, table.Row{"stx", "STX", fmt.Sprint(stxBalance), "", ""})

			for _, wallet := range wallets {
				for k, v := range wallet.FungibleTokens {
					balance, _ := strconv.ParseInt(v.Balance, 10, 64)
					if existingBalance, exists := fungibleTokenBalances[k]; exists {
						fungibleTokenBalances[k] = existingBalance + balance
					} else {
						fungibleTokenBalances[k] = balance
					}
				}
				for k, v := range wallet.NonFungibleTokens {
					count, _ := strconv.ParseInt(v.Count, 10, 64)
					if existingCount, exists := nonFungibleTokenCounts[k]; exists {
						nonFungibleTokenCounts[k] = existingCount + count
					} else {
						nonFungibleTokenCounts[k] = count
					}
				}
			}

			for k, balance := range fungibleTokenBalances {
				split := strings.Split(k, "::")
				contractName, err := props.HeroClient.GetContractReadOnly(split[0], "get-name", "string", []string{})
				if err != nil {
					rows = append(rows, table.Row{split[1], "Fungible", strconv.FormatInt(balance, 10), split[0], ""})
				}
				rows = append(rows, table.Row{split[1], "Fungible", strconv.FormatInt(balance, 10), split[0], strings.TrimSpace(contractName)})
			}

			for k, count := range nonFungibleTokenCounts {
				split := strings.Split(k, "::")
				rows = append(rows, table.Row{split[1], "Non-Fungible", strconv.FormatInt(count, 10), split[0], ""})
			}

			for _, row := range rows {
				for i, cell := range row {
					cellStr := fmt.Sprint(cell)
					if len(cellStr) > maxWidths[i] {
						maxWidths[i] = len(cellStr)
					}
				}
			}

			for i, maxWidth := range maxWidths {
				headers[i].Width = maxWidth
			}

			t := table.New(
				table.WithColumns(headers),
				table.WithRows(rows),
				table.WithFocused(true),
				table.WithStyles(common.TableStyles),
			)

			m := tableModel{table: t, client: props.HeroClient}
			if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
				props.Logger.Fatal().Err(err).Msg("Failed to run program")
			}

			return nil
		},
	}
}
