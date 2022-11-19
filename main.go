package main

import (
	"os"
	"time"

	"github.com/ItzAfroBoy/inv/fetch"
	"github.com/ItzAfroBoy/inv/helper"
	"github.com/ItzAfroBoy/inv/spinner"
	"github.com/ItzAfroBoy/inv/table"
	t "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var rows []t.Row

	sm := spinner.InitialModel()
	p := tea.NewProgram(sm)

	go func() {
		p.Send(fetch.ResMsg{Msg: "Fetching inventory", State: "running"})
		rows = fetch.Get(p)
		p.Send(fetch.ResMsg{Msg: "Saving inventory", State: "running"})
		time.Sleep(1 * time.Second)
		helper.Save(rows)
		p.Send(fetch.ResMsg{Msg: "Saving inventory", State: "complete"})
		time.Sleep(1 * time.Second)
		p.Send(fetch.ResMsg{Msg: "Loading inventory", State: "complete"})
		time.Sleep(500 * time.Millisecond)
		p.Quit()
	}()

	_, err := p.Run()
	helper.Check(err)

	if len(rows) == 0 {
		os.Exit(1)
	}

	tm := table.InitialModel(rows)
	_, err = tea.NewProgram(tm, tea.WithAltScreen()).Run()
	helper.Check(err)
}
