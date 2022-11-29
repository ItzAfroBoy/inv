package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ItzAfroBoy/inv/helper"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func Get(p *tea.Program, opts helper.Options) []table.Row {
	var rows []table.Row
	var res *http.Response

	if !opts.Cache {
		var err error
		res, err = http.Get("https://steamcommunity.com/inventory/76561198378367745/730/2")
		helper.Check(err)
	}

	if opts.Cache || res.StatusCode != 200 {
		p.Send(helper.ResMsg{Msg: "Fetching inventory", State: "failed"})
		time.Sleep(1 * time.Second)
		p.Send(helper.ResMsg{Msg: "Using cached inventory", State: "running"})

		for _, v := range helper.Open() {
			item := v.([]interface{})
			ix := item[0].(string)
			id := item[1].(string)
			name := item[2].(string)
			collection := item[3].(string)
			row := []string{ix, id, name, collection}

			if opts.Prices && opts.Cache {
				row = append(row, item[4].(string))
			}

			rows = append(rows, row)
		}

		p.Send(helper.ResMsg{Msg: "Using cached inventory", State: "complete"})
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
			}
		}
	}

	return rows
}

func GetPrices(slice []table.Row) []table.Row {
	var rows []table.Row

	for _, v := range slice {
		res, err := http.Get(fmt.Sprintf("https://steamcommunity.com/market/priceoverview?appid=730&currency=1&market_hash_name=%s", url.PathEscape(v[2])))
		helper.Check(err)

		if res.StatusCode != 200 {
			v = append(v, "$0.00")
			rows = append(rows, v)
		} else {
			var data map[string]interface{}

			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			helper.Check(err)

			err = json.Unmarshal(body, &data)
			helper.Check(err)

			if data["lowest_price"] != nil {
				v = append(v, data["lowest_price"].(string))
			} else {
				v = append(v, data["median_price"].(string))
			}

			rows = append(rows, v)
		}

		time.Sleep(3 * time.Second)
	}

	return rows
}

func GetInvValue(slice []table.Row) string {
	var value float64

	for _, v := range slice {
		x, _ := strconv.ParseFloat(v[4][1:], 64)
		value += x
	}

	return fmt.Sprintf("$%.2f", value)
}
