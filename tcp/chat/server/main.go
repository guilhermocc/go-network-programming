package main

import (
	"log"
	"net"
	"sync"
)

var addr = "localhost:8080"

func main() {
	server := NewChatServer()
	server.Run()
}

type ChatServer struct {
	clientsMtx sync.Mutex
	clients    map[string]net.Conn
	listener   net.Listener
}

func NewChatServer() *ChatServer {
	log.Printf("Starting server...")

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Listening on %s", addr)

	return &ChatServer{
		clientsMtx: sync.Mutex{},
		clients:    make(map[string]net.Conn),
		listener:   listener,
	}
}

func (cs *ChatServer) Run() {
	for {
		conn, err := cs.listener.Accept()
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		go cs.handleConnection(conn)
	}

}

func (cs *ChatServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	cs.clientsMtx.Lock()
	if _, ok := cs.clients[conn.RemoteAddr().String()]; ok {
		log.Printf("Client %v already has connection", conn.RemoteAddr())
		conn.Close()
		return
	}
	cs.clients[conn.RemoteAddr().String()] = conn
	cs.clientsMtx.Unlock()
	log.Printf("New connection from %v", conn.RemoteAddr())

	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				log.Printf("Client %v disconnected", conn.RemoteAddr())
			} else {
				log.Printf("Error reading from %v: %v", conn.RemoteAddr(), err)
			}
			conn.Close()
			return
		}
		cs.sendAll(conn.RemoteAddr(), string(buf))
	}
}

func (cs *ChatServer) sendAll(sender net.Addr, msg string) {
	cs.clientsMtx.Lock()
	for _, conn := range cs.clients {
		if conn.RemoteAddr().String() == sender.String() {
			continue
		}
		_, err := conn.Write([]byte(sender.String() + ": " + msg))
		if err != nil {
			log.Printf("Error writing to %v: %v", conn.RemoteAddr(), err)
			conn.Close()
			delete(cs.clients, conn.RemoteAddr().String())
		}
	}
	cs.clientsMtx.Unlock()
}
