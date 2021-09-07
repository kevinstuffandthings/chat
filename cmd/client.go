package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
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
			fmt.Printf("%s\n", string(buf[:l]))
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
	_, err = conn.Write([]byte(fmt.Sprintf("<Connect:@%s>", user)))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 4)
	_, err = conn.Read(buf)
	if err != nil || string(buf[:2]) != "OK" {
		return nil, errors.New(fmt.Sprintf("Connection improperly ack'd: <%s>", buf))
	}
	fmt.Println("You are connected!")
	return conn, nil
}
