package server

import (
	"fmt"
	"net"
)

type Message struct {
	Type     string
	Sender   string
	Contents string
}

func (m Message) Send(c net.Conn) {
	var text string
	if m.Type != "" {
		text = fmt.Sprintf("<%s> ", m.Type)
	}
	if m.Sender != "" {
		text += fmt.Sprintf("@%s", m.Sender)
	} else {
		text += "system"
	}
	c.Write([]byte(fmt.Sprintf("%s: %s", text, m.Contents)))
}
