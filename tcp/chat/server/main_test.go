package main

import (
	"bufio"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Buffer fixo
func Test1(t *testing.T) {

	// Processo do servidor

	// Cria socket e faz o Bind
	listener, err := net.Listen("tcp", "localhost:8080")
	require.NoError(t, err)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		t.Logf("New connection from %v", conn.RemoteAddr())
		require.NoError(t, err)

		buffer := make([]byte, 10)

		for {
			n, err := conn.Read(buffer)
			require.NoError(t, err)

			t.Logf("Read %d bytes", n)

			data := string(buffer[:n])

			t.Logf("Data: %s", data)

			n, err = conn.Write([]byte("Olá " + data))
			require.NoError(t, err)

			t.Logf("Write %d bytes", n)
		}

		conn.Close()
	}

}

// Using scanner
func Test2(t *testing.T) {
	// Processo do servidor

	// Cria socket e faz o Bind
	listener, err := net.Listen("tcp", "localhost:8080")
	require.NoError(t, err)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		t.Logf("New connection from %v", conn.RemoteAddr())
		require.NoError(t, err)

		scanner := bufio.NewScanner(conn)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
			_, err := conn.Write([]byte("Olá " + line + "\n"))
			require.NoError(t, err)
		}
		err = scanner.Err()
		if err != nil {
			t.Error(err)
		}

		conn.Close()
	}
}

// Error handling
func Test3(t *testing.T) {
	// Processo do servidor

	// Cria socket e faz o Bind
	listener, err := net.Listen("tcp", "localhost:8080")
	require.NoError(t, err)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		t.Logf("New connection from %v", conn.RemoteAddr())
		require.NoError(t, err)

		scanner := bufio.NewScanner(conn)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
			_, err := conn.Write([]byte("Olá " + line + "\n"))
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Println("timeout error:", err)
				time.Sleep(10 * time.Second)
				continue
			}
			t.Fatal("Unknown error when writing")
		}
		err = scanner.Err()
		if err != nil {
			t.Error(err)
		}

		conn.Close()
	}
}
