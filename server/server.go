package server

import (
	"fmt"
	"net"
	"sort"
	"sync"

	"github.com/kevinstuffandthings/chat/handshake"
)

type ChatServer struct {
	listener net.Listener
	users    map[string]net.Conn
	sync.Mutex
}

func New(listener net.Listener) ChatServer {
	return ChatServer{listener: listener}
}

func (s *ChatServer) Start(handlers ...Handler) {
	s.users = make(map[string]net.Conn)
	for {
		c, err := s.listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		u, err := handshake.Accept(c)
		if err != nil {
			fmt.Println("Unable to make handshake:", err)
			continue
		}

		s.connectUser(u, c)
		go func() {
			defer s.disconnectUser(u)

			buf := make([]byte, 1024)
			for {
				l, err := c.Read(buf)
				if err != nil {
					fmt.Println("User", u, "connection:", err)
					return
				}
				m := Message{Sender: u, Contents: string(buf[:l])}
				s.handleMessage(m, handlers)
			}
		}()
	}
}

func (s *ChatServer) OnlineUsers() []string {
	var users []string
	for u, l := range s.users {
		if l != nil {
			users = append(users, u)
		}
	}
	sort.Strings(users)
	return users
}

func (s *ChatServer) ConnectionFor(user string) net.Conn {
	c, ok := s.users[user]
	if !ok {
		return nil
	}
	return c
}

func (s *ChatServer) SendMessage(msg Message, conn net.Conn) error {
	_, err := conn.Write([]byte(msg.String()))
	return err
}

func (s *ChatServer) connectUser(user string, conn net.Conn) {
	s.Lock()
	defer s.Unlock()
	s.broadcast(Message{Contents: fmt.Sprintf("User @%s has entered the chat", user)})
	s.users[user] = conn
}

func (s *ChatServer) disconnectUser(user string) {
	s.Lock()
	defer s.Unlock()
	delete(s.users, user)
	s.broadcast(Message{Contents: fmt.Sprintf("User @%s has left the chat", user)})
}

func (s *ChatServer) broadcast(message Message) {
	for _, c := range s.users {
		s.SendMessage(message, c)
	}
}

func (s *ChatServer) handleMessage(message Message, handlers []Handler) {
	for _, h := range handlers {
		if cmd := h.Parse(message); cmd != nil {
			h.Execute(cmd, s)
			return
		}
	}
}
