package handshake

import (
	"errors"
	"fmt"
	"net"
	"regexp"
)

func Initiate(conn net.Conn, user string) error {
	_, err := conn.Write([]byte(fmt.Sprintf("<Connect:@%s>", user)))
	if err != nil {
		return err
	}

	buf := make([]byte, 4)
	_, err = conn.Read(buf)
	if err != nil || string(buf[:2]) != "OK" {
		return errors.New(fmt.Sprintf("Connection improperly ack'd: <%s>", buf))
	}

	return nil
}

func Accept(conn net.Conn) (string, error) {
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
	ack(conn, "OK")

	return uid, nil
}

func ack(conn net.Conn, msg string) error {
	if _, err := conn.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}
