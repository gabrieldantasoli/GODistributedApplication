package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

type HashQuery struct {
	Hash int `json:"hash"`
}

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
	case "download":
		if len(args) != 1 {
			fmt.Println("Uso: ./client download <hash>")
			return
		}
		hash := args[0]
		hashInt, err := strconv.Atoi(hash)
		if err != nil {
			log.Fatal("Erro ao converter hash:", err)
		}
		query := HashQuery{Hash: hashInt}
		sendJSON(conn, query)

		response := make([]byte, 4096)
		n, err := conn.Read(response)
		if err != nil {
			log.Fatal("Erro ao ler resposta do servidor:", err)
		}

		var ips []string
		err = json.Unmarshal(response[:n], &ips)
		fmt.Printf(ips)
		if err != nil {
			log.Fatal("Erro ao desserializar resposta:", err)
		}

		if len(ips) == 0 {
			fmt.Println("Nenhum host possui o arquivo com o hash especificado.")
			return
		}

		fmt.Printf("Tentando baixar o arquivo do IP: %s\n", ips[0])
		err = downloadFileFromHost(ips[0], hashInt, "dataset/")
		if err != nil {
			log.Fatalf("Erro ao baixar arquivo: %v", err)
		}
		fmt.Println("Download concluído com sucesso!")

	default:
		fmt.Println("Comando inválido. Use 'register', 'search' ou 'download'.")
		return
	}
}

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

func downloadFileFromHost(ip string, hash int, destFolder string) error {
	conn, err := net.Dial("tcp", ip+":9000")
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("Arquivo %s baixado com sucesso.\n")

	request := HashQuery{Hash: hash}
	err = json.NewEncoder(conn).Encode(request)
	if err != nil {
		return err
	}

	var fileInfo FileInfo
	err = json.NewDecoder(conn).Decode(&fileInfo)
	if err != nil {
		return err
	}

	destFile, err := os.Create(destFolder + fileInfo.FileName)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, conn)
	if err != nil {
		return err
	}

	return nil
}
