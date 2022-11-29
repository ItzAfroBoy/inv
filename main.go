package main

import (
	"flag"
	"fmt"
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
	var opts helper.Options

	prices := flag.Bool("prices", false, "Fetch item prices from Steam Market")
	cache := flag.Bool("cache", false, "Used cached results")

	sm := spinner.InitialModel()
	p := tea.NewProgram(sm)

	flag.Parse()

	opts = helper.Options{Prices: *prices, Cache: *cache}

	go func() {
		p.Send(helper.ResMsg{Msg: "Fetching inventory", State: "running"})
		rows = fetch.Get(p, opts)
		time.Sleep(1 * time.Second)
		p.Send(helper.ResMsg{Msg: "Fetching inventory", State: "complete"})
		p.Send(helper.ResMsg{Msg: fmt.Sprintf("%d items fetched", len(rows)), State: "complete"})
		time.Sleep(1 * time.Second)

		if opts.Prices && !opts.Cache {
			p.Send(helper.ResMsg{Msg: "Fetching prices", State: "running"})
			rows = fetch.GetPrices(rows)
			p.Send(helper.ResMsg{Msg: "Fetching prices", State: "complete"})
			time.Sleep(1 * time.Second)
			p.Send(helper.ResMsg{Msg: "Saving inventory", State: "running"})
			time.Sleep(1 * time.Second)
			helper.Save(rows)
			p.Send(helper.ResMsg{Msg: "Saving inventory", State: "complete"})
			time.Sleep(1 * time.Second)
		}

		p.Send(helper.ResMsg{Msg: "Loading inventory", State: "complete"})
		time.Sleep(500 * time.Millisecond)
		p.Quit()
	}()

	_, err := p.Run()
	helper.Check(err)

	if len(rows) == 0 {
		os.Exit(1)
	}

	tm := table.InitialModel(rows, opts)
	_, err = tea.NewProgram(tm, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
	helper.Check(err)
}
