package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ItzAfroBoy/inv/helper"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type ResMsg struct {
	Msg   string
	State string
}

func Get(p *tea.Program) []table.Row {
	var rows []table.Row

	res, err := http.Get("https://steamcommunity.com/inventory/76561198378367745/730/2")
	helper.Check(err)

	if res.StatusCode != 200 {
		p.Send(ResMsg{Msg: "Fetching inventory", State: "failed"})
		time.Sleep(1 * time.Second)
		p.Send(ResMsg{Msg: "Using cached inventory", State: "running"})

		for _, v := range helper.Open() {
			item := v.([]interface{})
			ix := item[0].(string)
			id := item[1].(string)
			name := item[2].(string)
			collection := item[3].(string)
			rows = append(rows, []string{ix, id, name, collection})
			time.Sleep(100 * time.Millisecond)
		}

		p.Send(ResMsg{Msg: "Using cached inventory", State: "complete"})
	} else {
		var data map[string]interface{}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		helper.Check(err)

		err = json.Unmarshal(body, &data)
		helper.Check(err)

		inv := data["descriptions"].([]interface{})
		assets := data["assets"].([]interface{})
		offset := 0

		for i, v := range assets {
			ix := fmt.Sprintf("%d", i+1)
			classid := v.(map[string]interface{})["classid"].(string)

			if classid == inv[i-offset].(map[string]interface{})["classid"] {
				item := inv[i-offset].(map[string]interface{})
				tags := item["tags"].([]interface{})
				id := v.(map[string]interface{})["assetid"].(string)
				name := item["market_hash_name"].(string)
				collection := tags[2].(map[string]interface{})["localized_tag_name"].(string)
				rows = append(rows, []string{ix, id, name, collection})
				time.Sleep(100 * time.Millisecond)
			} else {
				var name string
				var collection string
				id := v.(map[string]interface{})["assetid"].(string)
				for _, ev := range inv {
					item := ev.(map[string]interface{})
					if classid == item["classid"] {
						name = item["market_hash_name"].(string)
						collection = item["tags"].([]interface{})[2].(map[string]interface{})["localized_tag_name"].(string)
						break
					}
				}
				offset += 1
				rows = append(rows, []string{ix, id, name, collection})
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	time.Sleep(1 * time.Second)
	p.Send(ResMsg{Msg: "Fetching inventory", State: "complete"})
	p.Send(ResMsg{Msg: fmt.Sprintf("%d items fetched", len(rows)), State: "complete"})
	time.Sleep(1 * time.Second)
	return rows
}
