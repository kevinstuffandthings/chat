package server

import (
	"fmt"
	"net"
	"sort"
	"sync"

	"kevinstuffandthings/chat/handshake"
)

type ChatServer struct {
	listener net.Listener
	users    map[string]net.Conn
	mutex    sync.Mutex
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

func (s *ChatServer) connectUser(user string, conn net.Conn) {
	s.broadcast(Message{Contents: fmt.Sprintf("User @%s has entered the chat", user)})
	s.mutex.Lock()
	s.users[user] = conn
	s.mutex.Unlock()
}

func (s *ChatServer) disconnectUser(user string) {
	s.mutex.Lock()
	delete(s.users, user)
	s.mutex.Unlock()
	s.broadcast(Message{Contents: fmt.Sprintf("User @%s has left the chat", user)})
}

func (s *ChatServer) broadcast(message Message) {
	for _, c := range s.users {
		message.Send(c)
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
