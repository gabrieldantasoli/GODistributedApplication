package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

// Estrutura para consulta de hash
type HashQuery struct {
	Hash int `json:"hash"`
}

// Estrutura para informações de arquivo
type FileInfo struct {
	IP       string `json:"ip"`
	FileName string `json:"filename"`
	Hash     int    `json:"hash"`
	Action   string `json:"action"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: ./client <comando> <argumentos>")
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	switch command {
	case "register":
		// Envia comando de registro: register <IP> <hash1> <hash2> ...
		if len(args) < 2 {
			fmt.Println("Uso: ./client register <IP> <hash1> <hash2> ...")
			return
		}
		ip := args[0]
		hashes := args[1:]

		for _, hash := range hashes {
			hashInt, err := strconv.Atoi(hash)
			if err != nil {
				log.Fatal("Erro ao converter hash:", err)
			}
			fileInfo := FileInfo{
				IP:       ip,
				FileName: "somefile", // Substitua por um nome de arquivo adequado
				Hash:     hashInt,
				Action:   "add",
			}
			sendJSON(conn, fileInfo)
		}

	case "search":
		// Envia comando de busca: search <hash>
		if len(args) != 1 {
			fmt.Println("Uso: ./client search <hash>")
			return
		}
		hash := args[0]
		hashInt, err := strconv.Atoi(hash)
		if err != nil {
			log.Fatal("Erro ao converter hash:", err)
		}
		query := HashQuery{Hash: hashInt}
		sendJSON(conn, query)

	default:
		fmt.Println("Comando inválido. Use 'register' ou 'search'.")
		return
	}

	// Lê a resposta do servidor
	response := make([]byte, 4096)
	n, err := conn.Read(response)
	if err != nil {
		log.Fatal("Erro ao ler resposta do servidor:", err)
	}
	fmt.Println("Resposta do servidor:", string(response[:n]))
}

// Função para enviar dados como JSON
func sendJSON(conn net.Conn, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Erro ao serializar JSON:", err)
	}
	_, err = conn.Write(jsonData)
	if err != nil {
		log.Fatal("Erro ao enviar dados:", err)
	}
}
