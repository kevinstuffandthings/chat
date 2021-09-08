package server

type Command struct {
	Src  Message
	Name string
	Args string
}

type Handler interface {
	Parse(msg Message) *Command
	Execute(cmd *Command, srv *ChatServer) error
}

type BroadcastHandler struct{}

func (h BroadcastHandler) Parse(msg Message) *Command {
	return &Command{Src: msg, Args: msg.Contents}
}

func (h BroadcastHandler) Execute(cmd *Command, srv *ChatServer) error {
	for u, c := range srv.users {
		if u != cmd.Src.Sender {
			cmd.Src.Send(c)
		}
	}
	return nil
}
