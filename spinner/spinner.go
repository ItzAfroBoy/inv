package spinner

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ItzAfroBoy/inv/helper"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type model struct {
	spinner       spinner.Model
	actionSpinner spinner.Model
	actions       []string
	states        []string
	exit          bool
}

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Jump
	s.Style = lg.NewStyle().Foreground(lg.Color("#F15152"))

	sTwo := s
	sTwo.Spinner = spinner.Line
	sTwo.Style = s.Style.Copy().Foreground(lg.Color("#EDB183"))

	m := model{spinner: s, actionSpinner: sTwo}
	return m
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "ctrl+q":
			m.exit = true
			return m, tea.Sequence(tea.ClearScreen, tea.Quit)
		default:
			return m, nil
		}

	case helper.ResMsg:
		var state string

		switch msg.State {
		case "running":
			state = "spinner"
		case "complete":
			state = "✔️"
		case "failed":
			state = "❌"
		}

		i, found := sort.Find(len(m.actions), func(i int) int {
			return strings.Compare(msg.Msg, m.actions[i])
		})
		if !found {
			m.actions = append(m.actions, msg.Msg)
			m.states = append(m.states, state)
		} else {
			m.actions[i] = msg.Msg
			m.states[i] = state
		}

		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		m.actionSpinner, _ = m.actionSpinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	var str string

	if !m.exit {
		str = fmt.Sprintf("%s Running...", m.spinner.View())
	} else {
		str = fmt.Sprintf("%s Quitting...\n", m.spinner.View())
	}

	str += "\n\n"
	for i, v := range m.actions {
		if m.states[i] == "spinner" {
			str += fmt.Sprintf("%s %s\n", m.actionSpinner.View(), v)
		} else {
			str += fmt.Sprintf("%s %s\n", m.states[i], v)
		}
	}
	return str
}
