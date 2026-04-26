package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	SOCKET = "/tmp/go-cask-kv.sock"
)

const (
	OP_GET uint8 = 0
	OP_PUT uint8 = 1
	OP_DEL uint8 = 2
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Listen
	server, err := net.Listen("unix", SOCKET)
	if err != nil {
		return fmt.Errorf("Failed to listen: %s", err.Error())
	}
	defer server.Close()

	// run clean go routine
	go cleanup()

	fmt.Println("Listening on " + SOCKET)
	fmt.Println("Waiting for client...")

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %s", err.Error())
			continue
		}
		fmt.Println("client connected")
		go func() {
			defer conn.Close()
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered. Error:\n", r)
				}
			}()

			err := processConn(conn)
			if err != nil {
				_writeErr(conn, err.Error())
			}
		}()
	}
}

func cleanup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Remove(SOCKET)
	os.Exit(1)
}

func _read(conn net.Conn, length int) ([]byte, error) {
	buffer := make([]byte, length)
	rLen, err := conn.Read(buffer)
	if err != nil {
		return []byte{}, err
	}
	// Validate data length
	if rLen != length {
		return []byte{}, fmt.Errorf("Invalid size read: Got: %d, Expected: %d", rLen, length)
	}
	return buffer, nil
}

func _writeErr(conn net.Conn, msg string) error {
	_, err := fmt.Fprintf(conn, "ERROR: %s", msg)
	return err
}

func processConn(conn net.Conn) error {
	// Read op code
	opBuf, err := _read(conn, 1)
	if err != nil {
		return fmt.Errorf("Failed to read OPCODE data: %s", err.Error())
	}
	opCode := uint8(opBuf[0])

	// Read key size
	kBuf, err := _read(conn, 1)
	if err != nil {
		return fmt.Errorf("Failed to read key size: %s", err.Error())
	}
	keySize := uint8(kBuf[0])

	// Read key (k) bytes
	keyBuf, err := _read(conn, int(keySize))
	if err != nil {
		return fmt.Errorf("Failed to read key data: %s", err.Error())
	}

	switch opCode {
	case OP_GET:
		fmt.Printf("This is a GET operation to fetch key: %s\n", string(keyBuf))
	case OP_DEL:
		fmt.Printf("This is a DEL operation to delete key: %s\n", string(keyBuf))
	case OP_PUT:
		// 3. Read data size
		dBuf, err := _read(conn, 2)
		if err != nil {
			return fmt.Errorf("Failed to read data size: %s", err.Error())
		}
		dataSize := binary.BigEndian.Uint16(dBuf)

		// Read value (v) bytes
		valBuf, err := _read(conn, int(dataSize))
		if err != nil {
			return fmt.Errorf("Failed to read value data: %s", err.Error())
		}
		fmt.Printf("This is a PUT operation to insert key: %s and value: %s\n", string(keyBuf), string(valBuf))
	default:
		return fmt.Errorf("Invalid OPCODE")
	}

	_, err = conn.Write([]byte("Thanks! Got your message:"))
	return nil
}
