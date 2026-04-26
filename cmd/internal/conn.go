package internal

import "net"

const (
	SOCKET = "/tmp/go-cask-kv.sock"
)

func Connect() (net.Conn, error) {
	conn, err := net.Dial("unix", SOCKET)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
