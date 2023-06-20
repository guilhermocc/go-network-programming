package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	laddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1¨"), Port: 9091}
	raddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9090}
	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		log.Fatalf("Error dialing game server at %s: %s", raddr.String(), err)
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

	scanner := bufio.NewScanner(os.Stdin)

	scanner.Split(bufio.ScanLines)

	for {

		for scanner.Scan() {
			line := scanner.Text()
			_, err := conn.Write([]byte(line + "\n"))
			if err != nil {
				log.Fatalf("Error sending message to server: %v", err)
			}
		}

		if scanner.Err() != nil {
			fmt.Println("Error scanning input:", scanner.Err())
		}

	}
}
