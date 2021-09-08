package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
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

		u, err := handshake(c)
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

func handshake(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	l, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	rx := regexp.MustCompile(`^<Connect:@([A-Za-z0-9_-]+)>`)
	if !rx.Match(buf) {
		ack(conn, "ERR")
		return "", errors.New(fmt.Sprintf("Invalid connection header <%s>", buf))
	}
	uid := rx.ReplaceAllString(string(buf[:l]), "$1")
	fmt.Printf("User %s connected!\n", uid)
	ack(conn, "OK")

	return uid, nil
}

func ack(conn net.Conn, msg string) error {
	if _, err := conn.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}

func (s *ChatServer) online() []string {
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
	ts := ts()
	if message == "/users" {
		m := fmt.Sprintf("[%s] Online: %s", ts, strings.Join(s.online(), ", "))
		s.users[sender].Write([]byte(m))
	} else if message[0] == '@' {
		rx := regexp.MustCompile(`^@([^\s]+)\s+(.*)`)
		if match := rx.FindStringSubmatch(message); len(match) == 3 {
			if c, ok := s.users[match[1]]; ok {
				c.Write([]byte(fmt.Sprintf("[%s] @%s <private>: %s", ts, sender, match[2])))
			} else {
				s.users[sender].Write([]byte(fmt.Sprintf("[%s] User not online", ts)))
			}
		}
	} else {
		m := fmt.Sprintf("[%s] @%s: %s", ts, sender, message)
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
		c.Write([]byte(fmt.Sprintf("[%s] system: %s", ts(), message)))
	}
}

func ts() string {
	return time.Now().Format("15:04:05")
}
