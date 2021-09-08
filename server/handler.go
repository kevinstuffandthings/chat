package server

type Command struct {
	src  Message
	name string
	args string
}

func (c *Command) sender() string {
	return c.src.Sender
}

type Handler interface {
	Parse(msg Message) *Command
	Execute(cmd *Command, srv *ChatServer) error
}

type BroadcastHandler struct{}

func (h BroadcastHandler) Parse(msg Message) *Command {
	return &Command{src: msg, args: msg.Contents}
}

func (h BroadcastHandler) Execute(cmd *Command, srv *ChatServer) error {
	for u, c := range srv.users {
		if u != cmd.src.Sender {
			cmd.src.Send(c)
		}
	}
	return nil
}
