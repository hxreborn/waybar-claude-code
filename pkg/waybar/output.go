package waybar

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Output struct {
	Text       string `json:"text"`
	Tooltip    string `json:"tooltip,omitempty"`
	Class      string `json:"class,omitempty"`
	Percentage int    `json:"percentage"`
}

func (o *Output) Print() error {
	enc := json.NewEncoder(os.Stdout)
	return encode(enc, o)
}

func (o *Output) PrintTo(w io.Writer) error {
	enc := json.NewEncoder(w)
	return encode(enc, o)
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

func encode(enc *json.Encoder, o *Output) error {
	if err := enc.Encode(o); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	return nil
}
