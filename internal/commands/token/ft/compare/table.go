package compare

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
	"github.com/hashhavoc/teller/pkg/utils"
)

type tableModel struct {
	table          table.Model
	holdersTable   table.Model
	viewportBottom viewport.Model
	viewportTop    viewport.Model
	selected       table.Row

	client *hiro.APIClient
	logger log.Logger

	windowHeight int
	windowWidth  int

	sortAscending    bool
	lastSortedColumn int
	holdersView      bool
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
		default:
			switch msg.String() {
			case "esc":
				if m.table.Focused() {
					m.table.Blur()
				} else {
					m.table.Focus()
				}
			case "q", "ctrl+c":
				return m, tea.Quit
			case "1", "2", "3", "4", "5", "6", "7", "8", "9":
				columnIndex := int(msg.Runes[0] - '1')
				currentRows := m.table.Rows()
				columnCount := len(currentRows[0])

				if columnIndex < columnCount {
					if m.lastSortedColumn == columnIndex {
						m.sortAscending = !m.sortAscending
					} else {
						m.sortAscending = true
						m.lastSortedColumn = columnIndex
					}

					sort.SliceStable(currentRows, func(i, j int) bool {
						valI, errI := strconv.ParseFloat(currentRows[i][columnIndex], 64)
						valJ, errJ := strconv.ParseFloat(currentRows[j][columnIndex], 64)

						if errI == nil && errJ == nil {
							if m.sortAscending {
								return valI < valJ
							} else {
								return valI > valJ
							}
						}

						if m.sortAscending {
							return currentRows[i][columnIndex] < currentRows[j][columnIndex]
						} else {
							return currentRows[i][columnIndex] > currentRows[j][columnIndex]
						}
					})

					m.table.SetRows(currentRows)
				}
			case "enter":
				selectedRow := m.table.SelectedRow()
				contract, err := m.client.GetContractDetails(selectedRow[4])
				if err != nil {
					m.logger.Error().Err(err).Msg("Failed to get contract details")
					return m, nil
				}
				utils.OpenBrowser("https://explorer.hiro.so/txid/" + contract.TxID)
			case "a":
				err := common.WriteRowsToCSV(m.table.Rows(), "tokens.csv")
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
	if m.holdersView {
		view = m.holdersTable.View()
	} else {
		view = m.table.View()
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewportTop.View(),
		common.BaseTableStyle.Render(view),
		m.viewportBottom.View())
}
