package transactions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/urfave/cli/v2"
)

func CreateTransactionsCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "transactions",
		Usage: "Provides interactions with transactions",
		Subcommands: []*cli.Command{
			createSyncCommand(props),
		},
	}
}

func createSyncCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:      "sync",
		Usage:     "Sync transactions for a given principal",
		ArgsUsage: "principal",
		Action: func(c *cli.Context) error {
			principal := c.Args().First()
			if principal == "" {
				return cli.NewExitError("Please provide a principal (address or contract identifier)", 1)
			}
			return syncTransactions(props, principal)
		},
	}
}

func syncTransactions(props *props.AppProps, principal string) error {
	filename := fmt.Sprintf("%s_transactions.json", principal)

	// Load existing transactions from the file if it exists
	var existingTxs []Transaction
	if data, err := ioutil.ReadFile(filename); err == nil {
		if err := json.Unmarshal(data, &existingTxs); err != nil {
			return err
		}
	}

	// Fetch transactions from the Hiro API until the total matches the local count
	var allTxs []Transaction
	offset := 0
	limit := 50

	for {
		url := fmt.Sprintf("%s/extended/v1/address/%s/transactions?limit=%d&offset=%d", props.Config.Endpoints.Hiro, principal, limit, offset)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var txResp TransactionsResponse
		if err := json.NewDecoder(resp.Body).Decode(&txResp); err != nil {
			return err
		}

		allTxs = append(allTxs, txResp.Results...)

		if len(allTxs) >= txResp.Total {
			break
		}

		offset += limit
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

func mergeTxs(existingTxs, fetchedTxs []Transaction) []Transaction {
	txMap := make(map[string]bool)
	for _, tx := range existingTxs {
		txMap[tx.TxID] = true
	}

	var mergedTxs []Transaction
	mergedTxs = append(mergedTxs, existingTxs...)

	for _, tx := range fetchedTxs {
		if !txMap[tx.TxID] {
			mergedTxs = append(mergedTxs, tx)
			txMap[tx.TxID] = true
		}
	}

	return mergedTxs
}
