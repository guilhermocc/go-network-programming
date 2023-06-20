package main

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
)

const (
	ip   = "127.0.0.1"
	port = 9090
)

func main() {
	server := NewGameServer()
	server.Run()
}

type GameServer struct {
	guessNumberMtx sync.RWMutex
	guessNumber    int
	conn           *net.UDPConn
}

func NewGameServer() *GameServer {
	log.Printf("Iniciando servidor de jogo...")

	addr := &net.UDPAddr{IP: net.ParseIP(ip), Port: port}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Servidor de jogo iniciado em %v", addr)

	return &GameServer{
		conn:        conn,
		guessNumber: rand.Intn(100),
	}
}

func (cs *GameServer) Run() {
	log.Printf("Número a ser adivinhado: %d", cs.guessNumber)

	for {
		buffer := make([]byte, 4)
		n, clientAddr, err := cs.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		guess, err := strconv.Atoi(strings.ReplaceAll(string(buffer[:n]), "\n", ""))

		if err != nil {
			cs.conn.WriteToUDP([]byte("Número inválido\n"), clientAddr)
			continue
		}

		log.Printf("%s Tentou adivinhar %d", clientAddr, guess)

		cs.guessNumberMtx.RLock()
		guessed := guess == cs.guessNumber
		if guessed {
			cs.conn.WriteToUDP([]byte("Acertou\n"), clientAddr)
		} else if guess < cs.guessNumber {
			cs.conn.WriteToUDP([]byte("Maior\n"), clientAddr)
		} else {
			cs.conn.WriteToUDP([]byte("Menor\n"), clientAddr)
		}
		cs.guessNumberMtx.RUnlock()

		if guessed {
			cs.guessNumberMtx.Lock()
			cs.guessNumber = rand.Intn(100)
			cs.guessNumberMtx.Unlock()
			log.Printf("Número a ser adivinhado: %d", cs.guessNumber)
		}
	}
}
