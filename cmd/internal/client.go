package internal

import "fmt"

const (
	OP_GET uint8 = 0
	OP_PUT uint8 = 1
	OP_DEL uint8 = 2
)

func FormatGet(key string) []byte {
	// [OP_CODE][KEY SIZE][KEY]
	buffer := make([]byte, 0)
	buffer = append(buffer, OP_GET, uint8(len(key)))
	buffer = append(buffer, []byte(key)...)
	fmt.Printf("MESSAGE is: [%s]\n", string(buffer))
	return buffer
}

func FormatPut(key string, val string) []byte {
	// [OP_CODE][KEY SIZE][KEY][VALUE SIZE][VALUE]
	buffer := make([]byte, 0)
	buffer = append(buffer, OP_PUT, uint8(len(key)))
	buffer = append(buffer, []byte(key)...)
	buffer = append(buffer, uint8(len(val)))
	buffer = append(buffer, []byte(val)...)
	fmt.Printf("MESSAGE is: [%s]\n", string(buffer))
	return buffer
}

func FormatDel(key string) []byte {
	// [OP_CODE][KEY SIZE][KEY]
	buffer := make([]byte, 0)
	buffer = append(buffer, OP_DEL, uint8(len(key)))
	buffer = append(buffer, []byte(key)...)
	fmt.Printf("MESSAGE is: [%s]\n", string(buffer))
	return buffer
}
