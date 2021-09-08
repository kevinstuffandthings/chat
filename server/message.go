package server

import (
	"fmt"
	"net"
)

type Message struct {
	Sender   string
	Contents string
}

func (m Message) Send(c net.Conn) {
	var user string
	if m.Sender == "" {
		user = "system"
	} else {
		user = fmt.Sprintf("@%s", m.Sender)
	}
	c.Write([]byte(fmt.Sprintf("%s: %s", user, m.Contents)))
}
