package common

import (
	"encoding/csv"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var main = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#444444"}

var BaseTableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(main).
	Margin(1)

var TableSelectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("231")).
	Background(lipgloss.Color("033")).
	Bold(false)

var TableHeaderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(main).
	BorderBottom(true).
	Bold(true).
	Padding(0, 1)

var TableCellStyle = lipgloss.NewStyle().
	Padding(0, 1)

var TableStyles = table.Styles{
	Selected: TableSelectedStyle,
	Header:   TableHeaderStyle,
	Cell:     TableCellStyle,
}

const (
	TableHeightPadding = 10
	TableWidthPadding  = 5
)

// TableData represents a single row of data to be displayed in the table.
type TableData []string

// CreateTable creates a table with dynamic column widths based on the provided headers and data.
// headers: Titles of the columns.
// dataRows: Data to be displayed, where each TableData represents a row.
func CreateTable(headers []string, dataRows []TableData) table.Model {
	// Initialize table columns with default widths based on header titles
	columns := make([]table.Column, len(headers))
	for i, title := range headers {
		columns[i] = table.Column{Title: title, Width: len(title)}
	}

	// Prepare to track maximum widths
	maxWidths := make([]int, len(columns))
	for i, column := range columns {
		maxWidths[i] = column.Width
	}

	// Iterate over data to find maximum widths
	for _, row := range dataRows {
		for i, value := range row {
			if len(value) > maxWidths[i] {
				maxWidths[i] = len(value)
			}
		}
	}

	// Set the calculated widths to headers
	for i, maxWidth := range maxWidths {
		columns[i].Width = maxWidth
	}

	// Convert dataRows to table.Rows
	rows := make([]table.Row, len(dataRows))
	for i, dataRow := range dataRows {
		rows[i] = table.Row(dataRow)
	}

	// Create and return the table model
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithStyles(TableStyles),
	)

	return t
}

func InsertDecimal(str string, position int) string {
	// If position is 0, return the original string
	if position == 0 {
		return str
	}
	// Convert string to rune slice
	runes := []rune(str)

	// Calculate the position for inserting decimal
	insertPos := len(runes) - position

	// Check if insertPos is negative, indicating the need to prepend zeros
	if insertPos <= 0 {
		// Calculate the number of zeros to prepend
		numZeros := -(insertPos - 1) // Subtract 1 to account for the character at position 0
		// Prepend zeros
		for i := 0; i < numZeros; i++ {
			runes = append([]rune{'0'}, runes...)
		}
		insertPos = 1 // Adjust insert position to 1 as we've prepended zeros
	}

	// Insert decimal point
	runes = append(runes[:insertPos], append([]rune{'.'}, runes[insertPos:]...)...)

	// Convert back to string and return
	return string(runes)
}

func WriteRowsToCSV(rows []table.Row, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
