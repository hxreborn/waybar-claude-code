package waybar

import (
	"encoding/json"
	"fmt"
	"os"
)

type Output struct {
	Text       string      `json:"text"`
	Tooltip    string      `json:"tooltip,omitempty"`
	Class      interface{} `json:"class,omitempty"`
	Percentage int         `json:"percentage,omitempty"`
}

func (o *Output) Print() error {
	data, err := json.Marshal(o)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func Error(msg string) {
	o := &Output{
		Text:    "ERROR",
		Tooltip: msg,
		Class:   "error",
	}
	o.Print()
	os.Exit(1)
}
