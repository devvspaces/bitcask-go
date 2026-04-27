/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bitcask-go/cmd/internal"
	"fmt"
	"math"
	"os"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetch a value by the key",
	Run:   get,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func get(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("You need to provide a key")
		os.Exit(1)
	}
	key := args[0]
	if len(key) > math.MaxUint8 {
		fmt.Println("Key length can't be greater than 255 characters")
		os.Exit(1)
	}

	// establish connection
	conn, err := internal.Connect()
	if err != nil {
		fmt.Printf("Failed to establish connection: %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from error: %s\n", err.Error())
		}
	}()

	// Make a request
	_, err = conn.Write(internal.FormatGet(key))
	buffer := make([]byte, 1024)
	mLen, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	} else {
		// print response
		fmt.Println("Received: ", string(buffer[:mLen]))
	}
}
