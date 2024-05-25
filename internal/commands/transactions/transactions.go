package transactions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/internal/common"
	"github.com/hashhavoc/teller/pkg/api/hiro"
	"github.com/urfave/cli/v2"
)

func CreateTransactionsCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "transactions",
		Usage: "Provides interactions with transactions",
		Subcommands: []*cli.Command{
			createSyncCommand(props),
			createViewCommand(props),
		},
	}
}

func createViewCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "view",
		Usage: "View transactions for a given principal",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "principal",
				Aliases:  []string{"p"},
				Usage:    "Specify the principal",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			var rows []table.Row

			principal := c.String("principal")
			allTxs, err := props.HeroClient.GetTransactions(principal)
			if err != nil {
				return err
			}
			// Prepare the table
			// Prepare the table
			headers := []table.Column{
				{Title: "TxID", Width: len("0xce6a4bec9c1c3297e2a66cca212e3b29940b93066bedc4700931dea7e98c2d6a")},
				{Title: "Sender", Width: len("SP12BBFBGPH73KSM65QBF872GR6A0PGYR789R3HZG")},
				{Title: "Reciever", Width: len("SP12BBFBGPH73KSM65QBF872GR6A0PGYR789R3HZG")},
				{Title: "Status", Width: len("Status")},
				{Title: "Fee", Width: len("Fee")},
				{Title: "STX Sent", Width: len("STX Sent")},
				{Title: "STX Received", Width: len("STX Received")},
				{Title: "FT", Width: len("FT")},
				{Title: "NFT", Width: len("NFT")},
				{Title: "STX", Width: len("STX")},
			}

			maxWidths := make([]int, len(headers))
			for i, header := range headers {
				maxWidths[i] = header.Width
			}
			for _, tx := range allTxs {
				rows = append(rows, table.Row{
					tx.Tx.TxID,
					tx.Tx.SenderAddress,
					tx.Tx.TokenTransfer.RecipientAddress,
					tx.Tx.TxStatus,
					common.InsertDecimal(tx.Tx.FeeRate, 6),
					common.InsertDecimal(tx.Tx.TokenTransfer.Amount, 6),
					common.InsertDecimal(tx.StxReceived, 6),
					fmt.Sprint(tx.Events.Ft.Transfer),
					fmt.Sprint(tx.Events.Nft.Transfer),
					fmt.Sprint(tx.Events.Stx.Transfer),
				})
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
			m := tableModel{table: t, client: props.HeroClient, logger: props.Logger}
			if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
				props.Logger.Fatal().Err(err).Msg("Failed to run program")
			}
			return nil

		},
	}
}

func createSyncCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "sync",
		Usage: "Sync transactions for a given principal",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "principal",
				Aliases:  []string{"p"},
				Usage:    "Specify the principal",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			principal := c.String("principal")
			return syncTransactions(props, principal)
		},
	}
}

func syncTransactions(props *props.AppProps, principal string) error {
	filename := fmt.Sprintf("%s_transactions.json", principal)

	// Load existing transactions from the file if it exists
	var existingTxs []hiro.Transaction
	if data, err := ioutil.ReadFile(filename); err == nil {
		if err := json.Unmarshal(data, &existingTxs); err != nil {
			return err
		}
	}

	// Fetch transactions from the Hiro API until the total matches the local count
	var allTxs []hiro.Transaction
	allTxs, err := props.HeroClient.GetTransactions(principal)
	if err != nil {
		return err
	}

	// Merge fetched transactions with existing transactions
	mergedTxs := mergeTxs(existingTxs, allTxs)

	// Save updated transactions to the file
	data, err := json.MarshalIndent(mergedTxs, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	fmt.Printf("Synced %d transactions for principal %s\n", len(mergedTxs), principal)
	return nil
}

func mergeTxs(existingTxs, fetchedTxs []hiro.Transaction) []hiro.Transaction {
	txMap := make(map[string]bool)
	for _, tx := range existingTxs {
		txMap[tx.Tx.TxID] = true
	}

	var mergedTxs []hiro.Transaction
	mergedTxs = append(mergedTxs, existingTxs...)

	for _, tx := range fetchedTxs {
		if !txMap[tx.Tx.TxID] {
			mergedTxs = append(mergedTxs, tx)
			txMap[tx.Tx.TxID] = true
		}
	}

	return mergedTxs
}
