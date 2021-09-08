package handshake

import (
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
		return fmt.Errorf("connection improperly ack'd: <%s>", buf)
	}

	return nil
}

var handshakeRx *regexp.Regexp = regexp.MustCompile(`^<Connect:@([A-Za-z0-9_-]+)>`)

func Accept(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	l, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	if !handshakeRx.Match(buf) {
		ack(conn, "ERR")
		return "", fmt.Errorf("invalid connection header <%s>", buf)
	}
	uid := handshakeRx.ReplaceAllString(string(buf[:l]), "$1")
	ack(conn, "OK")

	return uid, nil
}

func ack(conn net.Conn, msg string) error {
	if _, err := conn.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}
