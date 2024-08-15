package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type model struct {
	table     table.Model
	user      string
	width     int
	getPrices bool
	sort      string
}

var baseStyle = lg.NewStyle().BorderStyle(lg.NormalBorder()).BorderForeground(lg.Color("240"))

func initialModel(length int) model {
	columns := []table.Column{
		{Title: "No", Width: len(parseIntString(length))},
		{Title: "Asset ID", Width: 11},
		{Title: "Item", Width: 20},
		{Title: "Collection", Width: 20},
		{Title: "Float", Width: 6},
		{Title: "Price", Width: 5},
	}

	t := table.New(table.WithColumns(columns), table.WithFocused(true))
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lg.ThickBorder()).BorderForeground(lg.Color("240")).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lg.Color("0")).Background(lg.Color("#D2D2D2")).Bold(false)

	t.SetStyles(s)
	return model{t, *user, len(parseIntString(length)), *getPrices, *sortTable}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
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

		}
	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			m.table.MoveUp(1)
		case tea.MouseButtonWheelDown:
			m.table.MoveDown(1)
		}
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 7)

	case addRowMsg:
		var rows []table.Row
		rows = append(rows, msg.rows...)

		if m.sort != "none" {
			switch m.sort {
			case "price":
				if m.getPrices {
					sort.Slice(rows, func(i, j int) bool {
						return parseFloat(rows[i][5][1:]) > parseFloat(rows[j][5][1:])
					})
				}
			case "float":
				sort.Slice(rows, func(i, j int) bool {
					return parseFloat(rows[i][4]) < parseFloat(rows[j][4])
				})
			case "collection":
				sort.Slice(rows, func(i, j int) bool {
					return rows[i][3] < rows[j][3]
				})
			case "item":
				sort.Slice(rows, func(i, j int) bool {
					return rows[i][2] < rows[j][2]
				})
			}

			for i, v := range rows {
				v[0] = fmt.Sprintf("%d", i+1)
			}
		}
		m.table.SetRows(rows)
		if m.table.Cursor() == (len(m.table.Rows()) - 2) {
			m.table.MoveDown(1)
		}
		return m, tea.SetWindowTitle(fmt.Sprintf("%s | %d items | %s", m.user, len(m.table.Rows()), parseInvValue(m.table.Rows())))
	case modifyLengthMsg:
		columns := []table.Column{
			{Title: "No", Width: m.width},
			{Title: "Asset ID", Width: 11},
			{Title: "Item", Width: msg.L1},
			{Title: "Collection", Width: msg.L2},
			{Title: "Float", Width: 6},
			{Title: "Price", Width: msg.L3},
		}
		m.table.SetColumns(columns)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if len(m.table.Rows()) > 0 {
		if strings.HasPrefix(m.table.SelectedRow()[2], "StatTrak") {
			s := table.DefaultStyles()
			s.Header = s.Header.BorderStyle(lg.ThickBorder()).BorderForeground(lg.Color("240")).BorderBottom(true).Bold(false)
			s.Selected = s.Selected.Foreground(lg.Color("0")).Background(lg.Color("#CF6A32")).Bold(false)
			m.table.SetStyles(s)
		} else if strings.HasPrefix(m.table.SelectedRow()[2], "Souvenir") {
			s := table.DefaultStyles()
			s.Header = s.Header.BorderStyle(lg.ThickBorder()).BorderForeground(lg.Color("240")).BorderBottom(true).Bold(false)
			s.Selected = s.Selected.Foreground(lg.Color("0")).Background(lg.Color("#FFD700")).Bold(false)
			m.table.SetStyles(s)
		}
	}

	return lg.JoinVertical(lg.Left, baseStyle.Render(m.table.View()), baseStyle.Render(lg.NewStyle().Render(fmt.Sprintf(" Items: %d | Value: %s ", len(m.table.Rows()), parseInvValue(m.table.Rows())))))
}
