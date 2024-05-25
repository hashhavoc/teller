package transactions

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/phuslu/log"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashhavoc/teller/internal/common"
	"github.com/hashhavoc/teller/pkg/api/hiro"
	"github.com/hashhavoc/teller/pkg/utils"
)

type tableModel struct {
	table            table.Model
	client           *hiro.APIClient
	sortAscending    bool
	lastSortedColumn int
	logger           log.Logger
}

func (m tableModel) Init() tea.Cmd {
	return nil
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - common.TableHeightPadding)
		return m, nil
	case tea.KeyMsg:
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
			selectedRow := m.table.SelectedRow()
			utils.OpenBrowser("https://explorer.hiro.so/txid/" + selectedRow[0])
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
		case "n":
			selectedRow := m.table.SelectedRow()
			m.table.SetRows(UpdateTableWithPrincipal(m, selectedRow[2]))

		case "b":
			selectedRow := m.table.SelectedRow()
			m.table.SetRows(UpdateTableWithPrincipal(m, selectedRow[1]))
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	return common.BaseTableStyle.Render(m.table.View())
}

func UpdateTableWithPrincipal(m tableModel, principal string) []table.Row {
	allTxs, err := m.client.GetTransactions(principal)
	if err == nil {
		var rows []table.Row
		for _, tx := range allTxs {
			rows = append(rows, table.Row{
				tx.Tx.TxID,
				common.ToName(tx.Tx.SenderAddress),
				common.ToName(tx.Tx.TokenTransfer.RecipientAddress),
				tx.Tx.TxStatus,
				common.InsertDecimal(tx.Tx.FeeRate, 6),
				common.InsertDecimal(tx.Tx.TokenTransfer.Amount, 6),
				common.InsertDecimal(tx.StxReceived, 6),
				fmt.Sprint(tx.Events.Ft.Transfer),
				fmt.Sprint(tx.Events.Nft.Transfer),
				fmt.Sprint(tx.Events.Stx.Transfer),
			})
		}
		return rows
	}
	return m.table.Rows()
}
