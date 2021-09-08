package main

import (
	"fmt"
	"net"
	"os"

	"kevinstuffandthings/chat/server"
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

	s := server.New(l)
	s.Start(server.SlashCmdHandler{}, server.BroadcastHandler{})
}
