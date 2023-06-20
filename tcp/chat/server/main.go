package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

const (
	ip   = "127.0.0.1"
	port = 8080
)

func main() {
	server := NewChatServer()
	server.Run()
}

type ChatServer struct {
	clientsMtx sync.Mutex
	clients    map[string]net.Conn
	listener   *net.TCPListener
}

func NewChatServer() *ChatServer {
	log.Printf("Iniciando servidor de chat...")

	addr := &net.TCPAddr{IP: net.ParseIP(ip), Port: port}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Servidor de chat iniciado em %v", addr)

	return &ChatServer{
		clientsMtx: sync.Mutex{},
		clients:    make(map[string]net.Conn),
		listener:   listener,
	}
}

func (cs *ChatServer) Run() {
	for {
		conn, err := cs.listener.AcceptTCP()
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		go cs.handleConnection(conn)
	}

}

func (cs *ChatServer) handleConnection(conn *net.TCPConn) {
	defer conn.Close()

	cs.createClient(conn)

	for {
		scanner := bufio.NewScanner(conn)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text() + "\n"
			cs.sendAll(conn.RemoteAddr(), line)
		}
		err := scanner.Err()
		if err != nil {
			log.Printf("Erro de leitura %v: %v", conn.RemoteAddr(), err)
		}

		cs.removeClient(conn)
		return
	}
}

func (cs *ChatServer) createClient(conn *net.TCPConn) {
	log.Printf("Cliente conectado %s", conn.RemoteAddr().String())
	cs.sendAll(conn.LocalAddr(), fmt.Sprintf("Cliente conectado %s\n", conn.RemoteAddr().String()))

	cs.clientsMtx.Lock()
	cs.clients[conn.RemoteAddr().String()] = conn
	cs.clientsMtx.Unlock()
}

func (cs *ChatServer) removeClient(conn *net.TCPConn) {
	cs.clientsMtx.Lock()
	delete(cs.clients, conn.RemoteAddr().String())
	cs.clientsMtx.Unlock()

	log.Printf("Cliente desconectado %s", conn.RemoteAddr().String())
	cs.sendAll(conn.LocalAddr(), fmt.Sprintf("Cliente desconectado %s\n", conn.RemoteAddr().String()))
}

func (cs *ChatServer) sendAll(sender net.Addr, msg string) {
	cs.clientsMtx.Lock()
	for _, conn := range cs.clients {
		if conn.RemoteAddr().String() == sender.String() {
			continue
		}
		_, err := conn.Write([]byte(sender.String() + ": " + msg))
		if err != nil {
			log.Printf("Erro de escrita %v: %v", conn.RemoteAddr(), err)
			delete(cs.clients, conn.RemoteAddr().String())
		}
	}
	cs.clientsMtx.Unlock()
}
