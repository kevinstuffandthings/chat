package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"kevinstuffandthings/chat/handshake"
)

func main() {
	port, user := os.Args[1], os.Args[2]
	conn, err := connect(fmt.Sprintf("localhost:%s", port), user)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go func() {
		buf := make([]byte, 1024)
		for {
			l, err := conn.Read(buf)
			if err != nil {
				panic(err)
			}
			fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), string(buf[:l]))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_, err := conn.Write([]byte(scanner.Text()))
		if err != nil {
			panic(err)
		}
	}
}

func connect(addr string, user string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	err = handshake.Initiate(conn, user)
	if err != nil {
		return nil, err
	}

	fmt.Println("You are connected!")
	return conn, nil
}
