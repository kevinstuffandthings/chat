package server

import (
	"fmt"
)

type Message struct {
	Type     string
	Sender   string
	Contents string
}

func (m Message) String() string {
	var text string
	if m.Type != "" {
		text = fmt.Sprintf("<%s> ", m.Type)
	}
	if m.Sender != "" {
		text += fmt.Sprintf("@%s", m.Sender)
	} else {
		text += "system"
	}
	return fmt.Sprintf("%s: %s", text, m.Contents)
}
