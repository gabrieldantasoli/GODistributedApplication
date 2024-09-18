package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"
)

type FileInfo struct {
	IP       string `json:"ip"`
	FileName string `json:"filename"`
	Hash     int    `json:"hash"`
	Action   string `json:"action"`
}

type HashQuery struct {
	Hash int `json:"hash"`
}

var fileMap = make(map[string]map[string]int)
var mutex sync.Mutex

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Servidor ouvindo na porta 8000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()

	var buf = make([]byte, 4096)

	n, err := c.Read(buf)
	if err != nil {
		if err != io.EOF {
			log.Println("Erro ao ler dados:", err)
		}
		return
	}

	log.Printf("Dados recebidos: %s\n", string(buf[:n]))

	var fileInfo FileInfo
	err = json.Unmarshal(buf[:n], &fileInfo)
	if err == nil && (fileInfo.Action == "add" || fileInfo.Action == "delete") {
		updateFileMap(fileInfo)
		_, err = c.Write([]byte("Arquivos e hashes atualizados com sucesso\n"))
		if err != nil {
			log.Println("Erro ao enviar resposta:", err)
		}
		return
	}

	var query HashQuery
	err = json.Unmarshal(buf[:n], &query)
	if err == nil && query.Hash != 0 {
		ips := getIPsForHash(query.Hash)
		response, err := json.Marshal(ips)
		if err != nil {
			log.Println("Erro ao serializar resposta:", err)
			return
		}
		_, err = c.Write(response)
		if err != nil {
			log.Println("Erro ao enviar resposta:", err)
		}
		return
	}

	log.Println("Formato de mensagem desconhecido")
}

func updateFileMap(fileInfo FileInfo) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := fileMap[fileInfo.IP]; !exists {
		fileMap[fileInfo.IP] = make(map[string]int)
	}

	if fileInfo.Action == "add" {
		fileMap[fileInfo.IP][fileInfo.FileName] = fileInfo.Hash
	} else if fileInfo.Action == "delete" {
		delete(fileMap[fileInfo.IP], fileInfo.FileName)
		if len(fileMap[fileInfo.IP]) == 0 {
			delete(fileMap, fileInfo.IP)
		}
	}

	log.Printf("Mapeamento atualizado para IP %s: %+v\n", fileInfo.IP, fileMap[fileInfo.IP])
}

func getIPsForHash(hash int) []string {
	mutex.Lock()
	defer mutex.Unlock()

	var ips []string
	for ip, files := range fileMap {
		for _, fileHash := range files {
			if fileHash == hash {
				ips = append(ips, ip)
				break
			}
		}
	}
	return ips
}
