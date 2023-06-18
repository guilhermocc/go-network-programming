package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

var serverAddr = "localhost:8080"

func main() {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Error dialing chat server at %s: %v", serverAddr, err)
	}

	go func() {
		for {
			buf := make([]byte, 1024)
			_, err := conn.Read(buf)
			if err != nil {
				log.Fatalf("Error reading from server: %v", err)
			}
			fmt.Print(string(buf))
		}
	}()

	// Create a new scanner to read from stdin
	scanner := bufio.NewScanner(os.Stdin)

	// Set the scanner to split on lines (default behavior)
	scanner.Split(bufio.ScanLines)

	for {

		// Read and process each line of the input
		for scanner.Scan() {
			line := scanner.Text()
			_, err := conn.Write([]byte(line + "\n"))
			if err != nil {
				log.Fatalf("Error sending message to server: %v", err)
			}
		}

		// Check for any errors encountered during scanning
		if scanner.Err() != nil {
			fmt.Println("Error scanning input:", scanner.Err())
		}

	}

}
