package contract

import (
	"fmt"
	"os"

	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/pkg/api/hiro"
	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
)

// CreateContractsCommand creates the contracts command and its subcommands.
func CreateContractsCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "contracts",
		Usage: "Provides interactions with contracts",
		Subcommands: []*cli.Command{
			createSourceCommand(props),
			createReadCommand(props),
			createViewCommand(props),
		},
	}
}

func createSourceCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "source",
		Usage: "view source of a contract",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "contract",
				Usage:    "Contract address",
				Aliases:  []string{"c"},
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			resp, err := props.HeroClient.GetContractSource(c.String("contract"))
			if err != nil {
				return err
			}
			fmt.Print(string(resp))
			return nil
		},
	}
}

func createReadCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "read",
		Usage: "execute read function of a contract",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "function",
				Usage:    "function to call to the contract",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "type",
				Usage:    "return type",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "contract",
				Usage:    "Contract address",
				Aliases:  []string{"c"},
				Required: true,
			},
		},
		ArgsUsage: "contract id",
		Action: func(c *cli.Context) error {
			resp, err := props.HeroClient.GetContractReadOnly(c.String("contract"), c.String("function"), c.String("type"), []string{})
			if err != nil {
				return err
			}
			fmt.Println(string(resp))
			return nil
		},
	}
}

func createViewCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:      "view",
		Usage:     "execute view of a contract",
		ArgsUsage: "contract id",
		Action: func(c *cli.Context) error {
			id := c.Args().First()
			if id == "" {
				return &ContractIDRequiredError{"contract id is required"}
			}
			resp := GetContractDetails(props.HeroClient, id)
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.SetStyle(table.StyleRounded)
			t.AppendHeader(table.Row{"Name", "Result"})
			for _, d := range resp {
				t.AppendRow([]interface{}{d.FunctionName, d.Result})
			}
			t.Render()
			return nil
		},
	}
}

func GetContractDetails(c *hiro.APIClient, id string) []ContractReadOnlyFunctionsSip10Response {
	functions := []ContractReadOnlyFunctionsSip10{}
	functions = append(functions, ContractReadOnlyFunctionsSip10{FunctionName: "get-name", ResponseType: "string"})
	functions = append(functions, ContractReadOnlyFunctionsSip10{FunctionName: "get-symbol", ResponseType: "string"})
	functions = append(functions, ContractReadOnlyFunctionsSip10{FunctionName: "get-token-uri", ResponseType: "string"})
	functions = append(functions, ContractReadOnlyFunctionsSip10{FunctionName: "get-decimals", ResponseType: "uint128"})
	functions = append(functions, ContractReadOnlyFunctionsSip10{FunctionName: "get-total-supply", ResponseType: "uint128"})
	details := make([]ContractReadOnlyFunctionsSip10Response, 0)
	for _, function := range functions {
		resp, _ := c.GetContractReadOnly(id, function.FunctionName, function.ResponseType, []string{})
		details = append(details, ContractReadOnlyFunctionsSip10Response{FunctionName: function.FunctionName, Result: resp})
	}

	return details
}
