package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Item struct {
	Text     string `json:"text"`
	Priority int    `json:"priority"`
	Position int    `json:"position"`
}

func SaveItems(filename string, items []Item) error {
	b, err := json.Marshal(items)
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	err = os.WriteFile(filename, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadItems(filename string) ([]Item, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return []Item{}, err
	}

	var items []Item
	if err := json.Unmarshal(b, &items); err != nil {
		return []Item{}, err
	}

	for i, _ := range items {
		items[i].Position = i + 1
	}

	return items, nil
}

func (i *Item) SetPriority(pri int) {
	switch pri {
	case 1:
		i.Priority = 1
	case 3:
		i.Priority = 3
	default:
		i.Priority = 2
	}

}

func (i *Item) PrettyPrint() string {
	if i.Priority == 1 {
		return "(1)"
	}
	if i.Priority == 3 {
		return "(3)"
	}
	return "  "
}

func (i *Item) Label() string {
	return strconv.Itoa(i.Position) + "."
}
