package table

import (
	"strings"

	"github.com/ItzAfroBoy/inv/fetch"
	"github.com/ItzAfroBoy/inv/helper"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type model struct {
	table table.Model
	len   int
	value string
}

var baseStyle = lg.NewStyle().BorderStyle(lg.DoubleBorder()).BorderForeground(lg.Color("240"))

func InitialModel(rows []table.Row, opts helper.Options) model {
	var value string
	m1, m2 := helper.Longest(rows)
	columns := []table.Column{
		{Title: "No", Width: 2},
		{Title: "Asset ID", Width: 11},
		{Title: "Item", Width: m1},
		{Title: "Collection", Width: m2},
	}

	if opts.Prices {
		columns = append(columns, table.Column{Title: "Price", Width: 5})
		value = fetch.GetInvValue(rows)
	}

	t := table.New(table.WithColumns(columns), table.WithRows(rows), table.WithFocused(true))
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lg.ThickBorder()).BorderForeground(lg.Color("240")).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lg.Color("0")).Background(lg.Color("#D2D2D2")).Bold(false)
	t.SetStyles(s)

	return model{t, len(rows), value}
}

func (m model) Init() tea.Cmd {
	return nil
}

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
	case tea.WindowSizeMsg:
		m.table.SetHeight(helper.Min(m.len, msg.Height-10))
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
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

	return baseStyle.Render(m.table.View()) + "\n" + baseStyle.Render(lg.NewStyle().Render(" Total Value: "+m.value+" ")) + "\n"
}
