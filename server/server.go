package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type FileRegistry struct {
	Files map[string][]string // Mapeia um hash para uma lista de IPs
	mu    sync.Mutex          // Para garantir que acessos concorrentes à tabela sejam seguros
}

func NewFileRegistry() *FileRegistry {
	return &FileRegistry{
		Files: make(map[string][]string),
	}
}

func (r *FileRegistry) Register(ip string, hashes []string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, hash := range hashes {
		r.Files[hash] = append(r.Files[hash], ip)
	}
}

func (r *FileRegistry) Search(hash string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.Files[hash]
}

func main() {
	registry := NewFileRegistry()

	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Servidor escutando na porta 8000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, registry)
	}
}

func handleConn(conn net.Conn, registry *FileRegistry) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) < 2 {
			fmt.Fprintln(conn, "Comando inválido")
			continue
		}

		command := parts[0]

		switch command {
		case "register":
			// Formato: register <IP> <hash1> <hash2> ...
			ip := parts[1]
			hashes := parts[2:]
			registry.Register(ip, hashes)
			fmt.Fprintln(conn, "Arquivos registrados com sucesso")

		case "search":
			// Formato: search <hash>
			hash := parts[1]
			ips := registry.Search(hash)
			if len(ips) > 0 {
				// Envia todos os IPs encontrados, separados por vírgulas
				fmt.Fprintln(conn, strings.Join(ips, ", "))
			} else {
				fmt.Fprintln(conn, "Nenhuma máquina encontrada para o hash")
			}

		default:
			fmt.Fprintln(conn, "Comando inválido")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Erro de leitura:", err)
	}
}
