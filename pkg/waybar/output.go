package waybar

import (
	"encoding/json"
	"fmt"
	"io"
)

type Output struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip,omitempty"`
	Class   string `json:"class,omitempty"`
}

func (o *Output) PrintTo(w io.Writer) error {
	if err := json.NewEncoder(w).Encode(o); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	return nil
}
