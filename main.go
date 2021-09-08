package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"

	"kevinstuffandthings/chat/handshake"
)

func main() {
	port := os.Args[1]
	addr := fmt.Sprintf("localhost:%s", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on", addr)

	server := ChatServer{listener: l}
	server.Start()
}

type ChatServer struct {
	listener net.Listener
	users    map[string]net.Conn
	mutex    sync.Mutex
}

func (s *ChatServer) Start() {
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

		s.broadcast(fmt.Sprintf("User @%s has entered the chat", u))
		s.mutex.Lock()
		s.users[u] = c
		s.handleMessage(u, "/users")
		go s.handleUserConnection(u, c)
		s.mutex.Unlock()
	}
}

func (s *ChatServer) onlineUsers() []string {
	var users []string
	for u, l := range s.users {
		if l != nil {
			users = append(users, u)
		}
	}
	sort.Strings(users)
	return users
}

func (s *ChatServer) handleUserConnection(user string, conn net.Conn) {
	defer func() {
		s.mutex.Lock()
		delete(s.users, user)
		s.mutex.Unlock()
		s.broadcast(fmt.Sprintf("User @%s has left the chat", user))
	}()

	buf := make([]byte, 1024)
	for {
		l, err := conn.Read(buf)
		if err != nil {
			fmt.Println("User", user, "connection:", err)
			return
		}

		s.handleMessage(user, string(buf[:l]))
	}
}

func (s *ChatServer) handleMessage(sender, message string) {
	if message == "/users" {
		m := fmt.Sprintf("Online: %s", strings.Join(s.onlineUsers(), ", "))
		s.users[sender].Write([]byte(m))
	} else if message[0] == '@' {
		rx := regexp.MustCompile(`^@([^\s]+)\s+(.*)`)
		if match := rx.FindStringSubmatch(message); len(match) == 3 {
			if c, ok := s.users[match[1]]; ok {
				c.Write([]byte(fmt.Sprintf("@%s <private>: %s", sender, match[2])))
			} else {
				s.users[sender].Write([]byte("User not online"))
			}
		}
	} else {
		m := fmt.Sprintf("@%s: %s", sender, message)
		fmt.Println(m)
		for u, c := range s.users {
			if u != sender {
				c.Write([]byte(m))
			}
		}
	}
}

func (s *ChatServer) broadcast(message string) {
	for _, c := range s.users {
		c.Write([]byte(fmt.Sprintf("system: %s", message)))
	}
}
