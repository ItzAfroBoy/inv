package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	queue := make(chan interface{})

	if *load {
		rows, l1, l2, l3 := importInv()
		m := initialModel(len(rows))
		p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion())

		go func() {
			p.Send(modifyLengthMsg{L1: l1, L2: l2, L3: l3})
			p.Send(addRowMsg{rows: rows})
		}()

		_, err := p.Run()
		check(err)
	} else {
		if *usecsf {
			go fetchCSFloatInv(queue)
		} else {
			go fetchInv(queue)
		}

		msg := <-queue
		length := msg.(int)
		m := initialModel(length)
		p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion())

		go func() {
			var rows []table.Row
			l1, l2, l3 := 0, 0, 0
			for elem := range queue {
				l1, l2, l3 = parseLength(elem.(table.Row), l1, l2, l3)
				rows = append(rows, elem.(table.Row))
				p.Send(modifyLengthMsg{L1: l1, L2: l2, L3: l3})
				p.Send(addRowMsg{rows: rows})
			}

			if *save {
				exportInv(rows)
			}
		}()

		_, err := p.Run()
		check(err)
	}
}
