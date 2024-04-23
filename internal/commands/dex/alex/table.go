package alex

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/phuslu/log"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashhavoc/teller/internal/common"
	"github.com/hashhavoc/teller/pkg/api/hiro"
)

type tableModel struct {
	table          table.Model
	holdersTable   table.Model
	viewportBottom viewport.Model
	viewportTop    viewport.Model
	client         *hiro.APIClient
	selected       table.Row
	logger         log.Logger

	windowHeight int
	windowWidth  int

	detailsView      bool
	holdersView      bool
	contractView     bool
	sortAscending    bool
	lastSortedColumn int
}

func (m tableModel) Init() tea.Cmd {
	m.viewportBottom.HighPerformanceRendering = true
	m.viewportTop.HighPerformanceRendering = true
	return tea.SetWindowTitle("Teller")
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		tcmd tea.Cmd
		bcmd tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowHeight = msg.Height
		m.windowWidth = msg.Width
		m.table.SetHeight(msg.Height - common.TableHeightPadding)
		m.holdersTable.SetHeight(msg.Height - common.TableHeightPadding)
		m.viewportBottom.Width = msg.Width
		m.viewportTop.Width = msg.Width
		return m, nil
	case tea.KeyMsg:
		// Use a switch case to handle different views
		switch {
		case m.holdersView:
			switch msg.String() {
			case "q":
				m.holdersView = false
				m.viewportTop.SetContent(fmt.Sprintf("Total: %s", fmt.Sprint(len(m.table.Rows()))))
				m.viewportBottom.SetContent("Press 'a' to export all addresses, 'h' to view holders")
				return m, nil
			case "a":
				filename := fmt.Sprintf("%s-holders.csv", m.selected[4])
				err := common.WriteRowsToCSV(m.holdersTable.Rows(), filename)
				if err != nil {
					m.logger.Error().Err(err).Msg("Failed to write rows to CSV file")
					return m, nil
				}
				m.viewportBottom.SetContent(fmt.Sprintf("Table dumped to %s", filename))
			}

			m.holdersTable, cmd = m.holdersTable.Update(msg)
			m.viewportTop, tcmd = m.viewportTop.Update(msg)
			m.viewportBottom, bcmd = m.viewportBottom.Update(msg)
			return m, tea.Batch(cmd, bcmd, tcmd)
		case m.contractView:
			switch msg.String() {
			case "q":
				// Switch back to details view from Contract ID view
				m.contractView = false
				return m, nil
			}
		case m.detailsView:
			switch msg.String() {
			case "c":
				// Switch to Contract ID view
				m.contractView = true
				return m, nil
			case "q":
				// Switch back to table view from details view
				m.detailsView = false
				return m, nil
			}
		default:
			// Handle the case when neither detailsView nor contractView is true
			switch msg.String() {
			case "esc":
				if m.table.Focused() {
					m.table.Blur()
				} else {
					m.table.Focus()
				}

			case "q", "ctrl+c":
				return m, tea.Quit
			case "enter":
				// Toggle to details view
				m.selected = m.table.SelectedRow()
				m.detailsView = true
			case "h":
				m.selected = m.table.SelectedRow()
				holders, err := m.client.GetTokenHolders(m.selected[4], 0)
				if err != nil {
					m.logger.Error().Err(err).Msg("Failed to get contract details")
					return m, nil
				}

				decimal, err := m.client.GetContractReadOnly(m.selected[4], "get-decimals", "uint128", []string{})
				if err != nil {
					m.logger.Error().Err(err).Msg("Failed to get contract details")
					return m, nil
				}
				i, err := strconv.Atoi(decimal)
				if err != nil {
					m.logger.Error().Err(err).Msg("Failed to convert string to integer")
					return m, nil
				}

				headers := []string{"Address", "Balance"}
				dataRows := generateHolderTableData(holders, i)
				t := common.CreateTable(headers, dataRows)

				m.holdersTable = t
				m.holdersView = true
				m.holdersTable.SetHeight(m.windowHeight - common.TableHeightPadding)
				m.viewportTop.SetContent(fmt.Sprintf("Total: %s", fmt.Sprint(len(dataRows))))
				m.viewportBottom.SetContent("Press 'a' to export all addresses")

			case "a":
				err := common.WriteRowsToCSV(m.table.Rows(), "alex-tokens.csv")
				if err != nil {
					return m, nil
				}
				m.logger.Info().Msg("Table dumped to tokens.csv")
			case "s":
				selectedRow := m.table.SelectedRow()
				contract, err := m.client.GetContractSource(selectedRow[4])
				if err != nil {
					m.logger.Error().Err(err).Msg("Failed to get contract source")
					return m, nil
				}

				var builder strings.Builder
				builder.WriteString(selectedRow[0])
				builder.WriteString("-")
				builder.WriteString(selectedRow[1])
				builder.WriteString("-")
				builder.WriteString(selectedRow[4])
				builder.WriteString(".clar")
				filename := builder.String()

				file, err := os.Create(filename)
				if err != nil {
					m.logger.Error().Err(err).Msg("Failed to create file")
					return m, nil
				}
				defer file.Close()

				_, err = file.WriteString(contract)
				if err != nil {
					m.logger.Error().Err(err).Msg("Failed to write to file")
					return m, nil
				}
				m.viewportBottom.SetContent(fmt.Sprintf("Contract source saved to %s", filename))
			case "1", "2", "3", "4", "5", "6", "7", "8", "9":
				columnIndex := int(msg.Runes[0] - '1') // Convert rune to int and adjust for 0-based indexing
				currentRows := m.table.Rows()
				columnCount := len(currentRows[0]) // Assuming all rows have the same number of columns

				if columnIndex < columnCount {
					// Check if the same column is being sorted again
					if m.lastSortedColumn == columnIndex {
						// Toggle the sorting direction
						m.sortAscending = !m.sortAscending
					} else {
						// New column, start with ascending order
						m.sortAscending = true
						m.lastSortedColumn = columnIndex
					}

					// Sort rows based on the dynamically chosen column and current sorting direction
					sort.SliceStable(currentRows, func(i, j int) bool {
						// Attempt to parse both values as int64
						valI, errI := strconv.ParseFloat(currentRows[i][columnIndex], 64)
						valJ, errJ := strconv.ParseFloat(currentRows[j][columnIndex], 64)

						// If both values are successfully parsed as int64, compare them as int64
						if errI == nil && errJ == nil {
							if m.sortAscending {
								return valI < valJ
							} else {
								return valI > valJ
							}
						}

						// If parsing into int64 fails for either value, compare them as strings
						if m.sortAscending {
							return currentRows[i][columnIndex] < currentRows[j][columnIndex]
						} else {
							return currentRows[i][columnIndex] > currentRows[j][columnIndex]
						}
					})

					// Set the sorted rows back to the table
					m.table.SetRows(currentRows)
				}
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	m.viewportTop, tcmd = m.viewportTop.Update(msg)
	m.viewportBottom, bcmd = m.viewportBottom.Update(msg)
	return m, tea.Batch(cmd, tcmd, bcmd)
}

func (m tableModel) View() string {
	var view string

	if m.contractView {
		// Render Contract ID view
		selectedRow := m.selected
		resp, _ := m.client.GetContractSource(selectedRow[4])
		return string(resp)
	} else if m.detailsView {
		// Render details view
		selectedRow := m.selected
		details := fmt.Sprintf("Selected Pair Details:\nBase: %s\nTarget: %s\nLast Price: %s\nLiquidity in USD: %s\nContract ID: %s",
			selectedRow[0], selectedRow[1], selectedRow[2], selectedRow[3], selectedRow[4])
		return details
	}
	if m.holdersView {
		view = m.holdersTable.View()
	} else {
		view = m.table.View()
	}
	// Render table view
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewportTop.View(),
		common.BaseTableStyle.Render(view),
		m.viewportBottom.View())
}

func generateHolderTableData(holders hiro.ContractHoldersResponse, i int) []common.TableData {
	var dataRows []common.TableData
	for k, d := range holders {
		strData := common.InsertDecimal(d, i)
		row := common.TableData{
			k,
			strData,
		}
		dataRows = append(dataRows, row)
	}
	return dataRows
}
