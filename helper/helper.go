package helper

import (
	"encoding/json"
	"fmt"
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
	User   string
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Open(user string) []interface{} {
	var dat []interface{}

	data, err := os.ReadFile(fmt.Sprintf("%s.json", user))
	Check(err)

	err = json.Unmarshal(data, &dat)
	Check(err)

	return dat
}

func Save(slice []table.Row, user string) {
	var rows []table.Row

	for _, v := range slice {
		rows = append(rows, v[1:])
	}

	data, _ := json.MarshalIndent(rows, "", "\t")
	err := os.WriteFile(fmt.Sprintf("%s.json", user), data, 0o644)
	Check(err)
}

func Longest(slice []table.Row) (L1 int, L2 int) {
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
