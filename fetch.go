package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
)

func fetchInv(ntfy chan interface{}) (items []table.Row, iconURLs [][]string) {
	var res *http.Response
	var err error
	var userID string
	var data map[string]interface{}

	res, err = http.PostForm("https://steamid.io/lookup", url.Values{"input": {*user}})
	check(err)
	userID = filepath.Base(res.Request.URL.Path)
	res, err = http.Get(fmt.Sprintf("https://steamcommunity.com/inventory/%s/730/2", userID))
	check(err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	check(err)

	check(json.Unmarshal(body, &data))
	inv := data["assets"].([]interface{})
	info := data["descriptions"].([]interface{})
	offset := 0

	ntfy <- len(inv)

	for i, v := range inv {
		var item table.Row
		ix := parseIntString(i + 1)
		classid := v.(map[string]interface{})["classid"].(string)
		assetid := v.(map[string]interface{})["assetid"].(string)

		if classid == info[i-offset].(map[string]interface{})["classid"] {
			item = parseItem(ix, assetid, userID, info[i-offset].(map[string]interface{}))
			items = append(items, item)
		} else {
			for _, _v := range info {
				if classid == _v.(map[string]interface{})["classid"] {
					item = parseItem(ix, assetid, userID, _v.(map[string]interface{}))
					items = append(items, item)
					offset += 1
					break
				}
			}
		}
		ntfy <- item
		if *getPrices {
			time.Sleep(3 * time.Second)
		}
	}

	if *getPrices {
		sendNotification(fmt.Sprintf("%d items @ %s", len(items), parseInvValue(items)))
	}
	close(ntfy)
	return
}

func fetchCSFloatInv(ntfy chan interface{}) (items []table.Row) {
	var data []interface{}

	req, err := http.NewRequest("GET", "https://csfloat.com/api/v1/me/inventory", nil)
	check(err)
	req.Header.Add("Authorization", *csfkey)
	res, err := http.DefaultClient.Do(req)
	check(err)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	check(err)
	check(json.Unmarshal(body, &data))

	ntfy <- len(data)

	for i, v := range data {
		ix := parseIntString(i + 1)
		item := parseCSFloatItem(ix, v.(map[string]interface{}))
		items = append(items, item)
		ntfy <- item
		if *getPrices {
			time.Sleep(3 * time.Second)
		}
	}
	if *getPrices && *notify {
		sendNotification(fmt.Sprintf("%d items @ %s", len(items), parseInvValue(items)))
	}
	close(ntfy)
	return
}

func fetchFloat(floatURL, userID, assetID string) string {
	var data map[string]interface{}

	floatURL = strings.ReplaceAll(floatURL, "%owner_steamid%", userID)
	floatURL = strings.ReplaceAll(floatURL, "%assetid%", assetID)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.csfloat.com/?url=%s", floatURL), nil)
	check(err)
	req.Header.Add("Origin", "https://csfloat.com/")
	res, err := http.DefaultClient.Do(req)
	check(err)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	check(err)
	check(json.Unmarshal(body, &data))
	float := data["iteminfo"].(map[string]interface{})["floatvalue"].(float64)

	return fmt.Sprintf("%f", float)[:6]
}

func fetchPrice(name string) string {
	res, err := http.Get(fmt.Sprintf("https://steamcommunity.com/market/priceoverview?appid=730&currency=1&market_hash_name=%s", url.PathEscape(name)))
	check(err)

	if res.StatusCode != 200 {
		return "$0.00"
	} else {
		var data map[string]interface{}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		check(err)
		check(json.Unmarshal(body, &data))

		if data["lowest_price"] != nil {
			return data["lowest_price"].(string)
		} else {
			return data["median_price"].(string)
		}
	}
}

func sendNotification(content string) {
	req, _ := http.NewRequest("POST", conf.NtfyEndpoint, strings.NewReader(content))
	req.Header.Set("Title", "CS:GO Inventory")
	_, err := http.DefaultClient.Do(req)
	check(err)
}
