package helper

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/bubbles/table"
)

type ResMsg struct {
	Msg   string
	State string
}

type Options struct {
	Prices bool
	Cache  bool
	Sort   string
	Order  string
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Open() []interface{} {
	var dat []interface{}

	data, err := os.ReadFile("inv.json")
	Check(err)

	err = json.Unmarshal(data, &dat)
	Check(err)

	return dat
}

func Save(slice []table.Row) {
	data, _ := json.MarshalIndent(slice, "", "\t")
	err := os.WriteFile("inv.json", data, 0644)
	Check(err)
}

func Longest(slice []table.Row) (int, int) {
	l1, l2 := 0, 0

	for _, v := range slice {
		func() {
			if len(v[2]) > l1 {
				l1 = len(v[2])
			}

			if len(v[3]) > l2 {
				l2 = len(v[3])
			}
		}()
	}

	return l1, l2
}

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
