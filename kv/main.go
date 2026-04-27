package main

import (
	"bitcask-go/kv/internal/hashmap"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	SOCKET = "/tmp/go-cask-kv.sock"
	DIR    = "./dir"
)

var (
	SEGMENT_FILE_PATH = filepath.Join(DIR, "segment.bin")
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

	// Ensure dir directory exists
	err = os.Mkdir(DIR, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Create db file
	f, err := os.OpenFile(SEGMENT_FILE_PATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("Failed to create segment file: %s", err.Error())
	}
	f.Close()

	// Create hash index
	hashmap.Init()

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
	key := string(keyBuf)

	switch opCode {
	case OP_GET:
		lod := hashmap.Get(key)

		// Open segment for reading
		f, err := os.Open(SEGMENT_FILE_PATH)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Seek(int64(lod.ValuePos), io.SeekStart)
		if err != nil {
			return err
		}

		buf := make([]byte, lod.ValueSize)
		_, err = io.ReadAtLeast(f, buf, int(lod.ValueSize))
		if err != nil {
			return err
		}

		_, err = conn.Write(buf)
		if err != nil {
			return err
		}
	case OP_DEL:
		fmt.Printf("This is a DEL operation to delete key: %s\n", key)
	case OP_PUT:
		// Read data size
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
		fmt.Printf("This is a PUT operation to insert key: %s and value: %s\n", key, string(valBuf))

		// Prepare data to write
		data := make([]byte, 0)
		tstamp := uint64(time.Now().Unix())
		btstamp := make([]byte, 8)
		binary.BigEndian.PutUint64(btstamp, tstamp)
		data = append(data, btstamp...)
		data = append(data, keySize)
		data = append(data, dBuf...)
		data = append(data, keyBuf...)
		data = append(data, valBuf...)

		// Open segment for writing
		f, err := os.OpenFile(SEGMENT_FILE_PATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Write(data)
		if err != nil {
			return err
		}
		err = f.Sync()
		if err != nil {
			return err
		}
		endPos, _ := f.Seek(0, io.SeekCurrent)
		valuePos := endPos - int64(dataSize)

		// Now we need to update the hash map
		hashmap.Upsert(key, hashmap.Dir{
			FileId:    SEGMENT_FILE_PATH,
			Timestamp: tstamp,
			ValueSize: dataSize,
			ValuePos:  uint64(valuePos),
		})

		_, err = conn.Write([]byte("Success!"))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Invalid OPCODE")
	}

	return nil
}
