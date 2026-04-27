package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	err := os.Mkdir("./dir", 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	path := filepath.Join("./dir", "db")

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.Write([]byte{12, 34, 56, 21})
	if err != nil {
		panic(err)
	}
	err = w.Flush()
	if err != nil {
		panic(err)
	}

	f2, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f2.Close()

	buf := make([]byte, 2)
	offset, err := f2.Seek(1, io.SeekStart)
	if err != nil {
		panic(err)
	}
	fmt.Println("File offset at ", offset)

	_, err = io.ReadAtLeast(f2, buf, 2)
	if err != nil {
		panic(err)
	}

	fmt.Println("Data read back:", buf)
}
