package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: ./bodocongo_bay <comando> <argumentos>")
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
			fmt.Println("Uso: ./bodocongo_bay register <IP> <hash1> <hash2> ...")
			return
		}
		ip := args[0]
		hashes := args[1:]
		fmt.Fprintf(conn, "register %s %s\n", ip, strings.Join(hashes, " "))

	case "search":
		// Envia comando de busca: search <hash>
		if len(args) != 1 {
			fmt.Println("Uso: ./bodocongo_bay search <hash>")
			return
		}
		hash := args[0]
		fmt.Fprintf(conn, "search %s\n", hash)

	default:
		fmt.Println("Comando inválido. Use 'register' ou 'search'.")
		return
	}

	// Lê a resposta do servidor
	response, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print(response)
}
