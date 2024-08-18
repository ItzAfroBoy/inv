package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"gopkg.in/yaml.v3"
)

type ITEM struct {
	AssetID    string
	Name       string
	Collection string
	Float      string
	Price      string
}

type CONFIG struct {
	User       string
	SteamID    string
	CSFloatKey string
	UseCSFloat bool
	NtfyEndpoint string
}

type addRowMsg struct {
	rows []table.Row
}

type modifyLengthMsg struct {
	L1 int
	L2 int
	L3 int
}

var conf CONFIG
var getPrices *bool
var printConfig *bool
var sortTable *string
var user *string
var save *bool
var load *bool
var usecsf *bool
var csfkey *string
var notify *bool

func init() {
	getPrices = flag.Bool("prices", false, "Fetch item prices from Steam Market")
	printConfig = flag.Bool("print-config", false, "Print config file to command line")
	sortTable = flag.String("sort", "none", "Sort by price, item, collection or float")
	user = flag.String("user", "none", "Steam user's inventory to fetch")
	csfkey = flag.String("csf-key", "none", "CSFloat API key")
	usecsf = flag.Bool("use-csf", false, "Use CSFloat data")
	save = flag.Bool("export", false, "Save your inventory to a JSON file")
	load = flag.Bool("import", false, "Import your inventory from a JSON file")
	notify = flag.Bool("notify", false, "Send a notification to your phone (url must be set in conf file)")

	flag.Parse()

	fp := filepath.Join(os.Getenv("HOME"), ".config", "inv", "config.yaml")
	_, err := os.Stat(fp)
	if err == nil {
		body, err := os.ReadFile(fp)
		check(err)
		check(yaml.Unmarshal(body, &conf))
	}

	checkFlags()
	if *printConfig {
		printConfigFile()
	}
}

func printConfigFile() {
	fmt.Println("User: ", conf.User)
	fmt.Println("CSFloat Key: ", conf.CSFloatKey)
	fmt.Println("Use CSFloat API: ", conf.UseCSFloat)
	fmt.Println("SteamID: ", conf.SteamID)
	fmt.Println("Ntfy Endpoint: ", conf.NtfyEndpoint)
	os.Exit(0)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checkFlags() {
	if *user == "none" {
		if conf.User == "" {
			os.Exit(1)
		}
		*user = conf.User
	}
	if !*usecsf {
		if conf.UseCSFloat {
			*usecsf = true
			if *csfkey == "none" {
				if conf.CSFloatKey == "" {
					os.Exit(1)
				}
				*csfkey = conf.CSFloatKey
			}
		}
	} else if *csfkey == "none" {
		if conf.CSFloatKey == "" {
			os.Exit(1)
		}
		*csfkey = conf.CSFloatKey
	}
}

func parseFloat(str string) (float float64) {
	float, _ = strconv.ParseFloat(str, 64)
	return
}

func parseIntString(num int) (str string) {
	str = fmt.Sprintf("%d", num)
	return
}

func parseLength(row table.Row, l1, l2, l3 int) (int, int, int) {
	if len(row[2]) > l1 {
		l1 = len(row[2])
	}
	if len(row[3]) > l2 {
		l2 = len(row[3])
	}
	if len(row[5]) > l3 {
		l3 = len(row[5])
	}
	return l1, l2, l3
}

func parseItem(ix, assetid, userID string, item map[string]interface{}) table.Row {
	var price string
	tags := item["tags"].([]interface{})
	name := item["market_hash_name"].(string)
	collection := tags[2].(map[string]interface{})["localized_tag_name"].(string)
	float := fetchFloat(item["actions"].([]interface{})[0].(map[string]interface{})["link"].(string), userID, assetid)
	if *getPrices {
		price = fetchPrice(name)
	} else {
		price = "$0.00"
	}
	return table.Row{ix, assetid, name, collection, float, price}
}

func parseCSFloatItem(ix string, item map[string]interface{}) table.Row {
	var price string
	assetid := item["asset_id"].(string)
	name := item["market_hash_name"].(string)
	collection := item["collection"].(string)
	float := fmt.Sprintf("%f", item["float_value"].(float64))[:6]
	if *getPrices {
		price = fetchPrice(name)
	} else {
		price = "$0.00"
	}
	return table.Row{ix, assetid, name, collection, float, price}
}

func parseInvValue(items []table.Row) string {
	var value float64
	for _, v := range items {
		price := parseFloat(v[5][1:])
		value += price
	}
	return fmt.Sprintf("$%.2f", value)
}

func exportInv(items []table.Row) {
	var out []map[string]string

	for _, v := range items {
		item := map[string]string{
			"assetid":    v[1],
			"name":       v[2],
			"collection": v[3],
			"float":      v[4],
			"price":      v[5],
		}

		out = append(out, item)
	}

	data, err := json.MarshalIndent(out, "", "\t")
	check(err)
	err = os.WriteFile(fmt.Sprintf("%s.json", *user), data, 0o644)
	check(err)
}

func importInv() ([]table.Row, int, int, int) {
	var inv []ITEM

	body, err := os.ReadFile(fmt.Sprintf("%s.json", *user))
	check(err)
	err = json.Unmarshal(body, &inv)
	check(err)

	var rows []table.Row
	l1, l2, l3 := 0, 0, 0

	for i, v := range inv {
		ix := fmt.Sprintf("%d", i+1)
		row := table.Row{ix, v.AssetID, v.Name, v.Collection, v.Float, v.Price}
		l1, l2, l3 = parseLength(row, l1, l2, l3)
		rows = append(rows, row)
	}

	return rows, l1, l2, l3
}
