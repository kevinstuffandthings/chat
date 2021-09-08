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

	cmd := &Command{Src: msg}
	rx := regexp.MustCompile(`^\/([^\s]+)(.*)`)
	match := rx.FindStringSubmatch(msg.Contents)

	if len(match) >= 1 {
		cmd.Name = match[1]
	}
	if len(match) >= 2 {
		cmd.Args = strings.TrimSpace(match[2])
	}
	return cmd
}

func (h SlashCmdHandler) Execute(cmd *Command, srv *ChatServer) error {
	user := srv.ConnectionFor(cmd.Src.Sender)
	switch cmd.Name {
	case "users":
		msg := Message{Contents: fmt.Sprintf("Online: %s", strings.Join(srv.OnlineUsers(), ", "))}
		msg.Send(user)
	case "dm":
		rx := regexp.MustCompile(`^@([^\s]+)\s+(.*)`)
		if match := rx.FindStringSubmatch(cmd.Args); len(match) == 3 {
			if c := srv.ConnectionFor(match[1]); c != nil {
				Message{Type: "DM", Sender: cmd.Src.Sender, Contents: match[2]}.Send(c)
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
