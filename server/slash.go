package server

import (
	"fmt"
	"regexp"
	"strings"
)

type SlashCmdHandler struct{}

func (h SlashCmdHandler) Parse(msg Message) *Command {
	if msg.Contents[0] != '/' {
		return nil
	}

	cmd := &Command{src: msg}
	rx := regexp.MustCompile(`^\/([^\s]+)(.*)`)
	match := rx.FindStringSubmatch(msg.Contents)

	if len(match) >= 1 {
		cmd.name = match[1]
	}
	if len(match) >= 2 {
		cmd.args = strings.TrimSpace(match[2])
	}
	return cmd
}

func (h SlashCmdHandler) Execute(cmd *Command, server *ChatServer) error {
	user := server.users[cmd.sender()]
	switch cmd.name {
	case "users":
		msg := Message{Contents: fmt.Sprintf("Online: %s", strings.Join(server.OnlineUsers(), ", "))}
		msg.Send(server.users[cmd.sender()])
	case "dm":
		rx := regexp.MustCompile(`^@([^\s]+)\s+(.*)`)
		if match := rx.FindStringSubmatch(cmd.args); len(match) == 3 {
			if c, ok := server.users[match[1]]; ok {
				c.Write([]byte(fmt.Sprintf("@%s <private>: %s", cmd.sender(), match[2])))
			} else {
				Message{Contents: "User not online"}.Send(user)
			}
		} else {
			Message{Contents: "Unable to understand dm"}.Send(user)
		}
	case "help":
		Message{Contents: "/users: get a list of all online users\n/dm @USER MSG: send a direct message to the specified user"}.Send(user)
	default:
		Message{Contents: "Unknown command"}.Send(user)
	}
	return nil
}
