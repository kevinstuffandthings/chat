package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
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

	server := ChatServer{Listener: l}
	server.Start()
}

type ChatServer struct {
	Listener net.Listener
	users    map[string]net.Conn
}

func (s *ChatServer) Start() {
	s.users = make(map[string]net.Conn)
	for {
		c, err := s.Listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		u, err := handshake(c)
		if err != nil {
			fmt.Println("Unable to make handshake:", err)
			continue
		}
		s.users[u] = c
		go s.handleUserConnection(u, c)
	}
}

func handshake(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	l, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	rx := regexp.MustCompile("^<Connect:@([A-Za-z0-9_-]+)>")
	if !rx.Match(buf) {
		ack(conn, "ERR")
		return "", errors.New(fmt.Sprintf("Invalid connection header <%s>", buf))
	}
	uid := rx.ReplaceAllString(string(buf[:l]), "$1")
	fmt.Printf("User <%s> connected!\n", uid)
	ack(conn, "OK")

	return uid, nil
}

func ack(conn net.Conn, msg string) error {
	if _, err := conn.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}

func (s *ChatServer) handleUserConnection(user string, conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		l, err := conn.Read(buf)
		if err != nil {
			fmt.Println("User", user, "connection:", err)
			return
		}

		msg := fmt.Sprintf("@%s: %s", user, string(buf[:l]))
		fmt.Println(msg)
		for u, c := range s.users {
			if u != user {
				c.Write([]byte(msg))
			}
		}
	}
}
