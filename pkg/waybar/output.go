package waybar

import (
	"encoding/json"
	"fmt"
	"os"
)

type Output struct {
	Text       string `json:"text"`
	Tooltip    string `json:"tooltip,omitempty"`
	Class      string `json:"class,omitempty"`
	Percentage int    `json:"percentage,omitempty"`
}

func (o *Output) Print() error {
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(o); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	return nil
}

func Error(msg string) {
	o := &Output{
		Text:    "ERROR",
		Tooltip: msg,
		Class:   "error",
	}
	_ = o.Print()
	os.Exit(1)
}
