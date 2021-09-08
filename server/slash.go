package server

import (
	"fmt"
	"regexp"
	"strings"
)

type SlashCmdHandler struct{}

var cmdRx *regexp.Regexp = regexp.MustCompile(`^\/([^\s]+)(.*)`)

func (h SlashCmdHandler) Parse(msg Message) *Command {
	if msg.Contents[0] != '/' {
		return nil
	}

	cmd := &Command{Src: msg}
	match := cmdRx.FindStringSubmatch(msg.Contents)

	if len(match) >= 1 {
		cmd.Name = match[1]
	}
	if len(match) >= 2 {
		cmd.Args = strings.TrimSpace(match[2])
	}
	return cmd
}

var (
	dmRx *regexp.Regexp = regexp.MustCompile(`^@([^\s]+)\s+(.*)`)

	msgUserNotOnline Message = Message{Contents: "User not online"}
	msgBadDm         Message = Message{Contents: "Unable to understand dm"}
	msgUnknownCmd    Message = Message{Contents: "Unknown command"}
	msgHelp          Message = Message{Contents: "/users: get a list of all online users\n/dm @USER MSG: send a direct message to the specified user"}
)

func (h SlashCmdHandler) Execute(cmd *Command, srv *ChatServer) error {
	user := srv.ConnectionFor(cmd.Src.Sender)
	switch cmd.Name {
	case "users":
		msg := Message{Contents: fmt.Sprintf("Online: %s", strings.Join(srv.OnlineUsers(), ", "))}
		srv.SendMessage(msg, user)
	case "dm":
		if match := dmRx.FindStringSubmatch(cmd.Args); len(match) == 3 {
			if c := srv.ConnectionFor(match[1]); c != nil {
				srv.SendMessage(Message{Type: "DM", Sender: cmd.Src.Sender, Contents: match[2]}, c)
			} else {
				srv.SendMessage(msgUserNotOnline, user)
			}
		} else {
			srv.SendMessage(msgBadDm, user)
		}
	case "help":
		srv.SendMessage(msgHelp, user)
	default:
		srv.SendMessage(msgUnknownCmd, user)
	}
	return nil
}
