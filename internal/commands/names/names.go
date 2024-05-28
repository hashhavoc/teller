package names

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/hashhavoc/teller/internal/common"
	"github.com/jszwec/csvutil"
	"github.com/urfave/cli/v2"
)

func CreateNameCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "names",
		Usage: "Provides interactions with names",
		Subcommands: []*cli.Command{
			createViewCommand(props),
			createLookupCommand(props),
			createSyncCommand(props),
		},
	}
}

func createLookupCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "lookup",
		Usage: "lookup names",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "The name to lookup",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			var rows []table.Row

			theName, err := props.HeroClient.GetName(c.String("name"))
			if err != nil {
				return err
			}

			headers := []table.Column{
				{Title: "Name", Width: len("0xce6a4bec9c1c3297e2a66cca212e3b29940b93066bedc4700931dea7e98c2d6a")},
				{Title: "Address", Width: len("SPSR9XHHRG3XYQ59A13Z1WSWESPRDBXCGX9VXEMP")},
				{Title: "Expire Block", Width: len("5000000")},
				{Title: "Registered Block", Width: len("5000000")},
			}

			maxWidths := make([]int, len(headers))
			for i, header := range headers {
				maxWidths[i] = header.Width
			}
			rows = append(rows, table.Row{
				c.String("name"),
				theName.Address,
			})

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

			m := tableModel{table: t, client: props.HeroClient, logger: props.Logger, page: 1}
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
		Usage: "sync names",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "the filename to output to",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "type",
				Aliases:  []string{"t"},
				Usage:    "The format to output to (json, csv)",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			return syncNames(props, c.String("type"), c.String("file"))
		},
	}
}

func syncNames(props *props.AppProps, format string, filename string) error {
	var data []byte
	var err error

	allNames, err := props.HeroClient.GetAllNames()
	if err != nil {
		return err
	}

	switch format {
	case "csv":
		data, err = csvutil.Marshal(allNames)
	case "json":
		data, err = json.MarshalIndent(allNames, "", "  ")
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	fmt.Printf("Synced %d BNS Name Records to %s\n", len(allNames), filename)
	return nil
}

func createViewCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "view",
		Usage: "View names",
		Action: func(c *cli.Context) error {
			var rows []table.Row

			allNames, err := props.HeroClient.GetAllNames()
			if err != nil {
				return err
			}

			// Prepare the table
			headers := []table.Column{
				{Title: "Name", Width: len("0xce6a4bec9c1c3297e2a66cca212e3b29940b93066bedc4700931dea7e98c2d6a")},
				{Title: "Address", Width: len("SPSR9XHHRG3XYQ59A13Z1WSWESPRDBXCGX9VXEMP")},
				{Title: "Expire Block", Width: len("500000000000000")},
				{Title: "Registered Block", Width: len("500000000000000")},
			}

			maxWidths := make([]int, len(headers))
			for i, header := range headers {
				maxWidths[i] = header.Width
			}
			for _, name := range allNames {
				rows = append(rows, table.Row{
					name.Name,
					name.Address,
					fmt.Sprintf("%d", name.ExpireBlock),
					fmt.Sprintf("%d", name.RegisteredAt),
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
			m := tableModel{table: t, client: props.HeroClient, logger: props.Logger, page: 1}
			if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
				props.Logger.Fatal().Err(err).Msg("Failed to run program")
			}
			return nil

		},
	}
}
