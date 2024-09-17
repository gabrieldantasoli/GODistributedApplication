package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"
)

// Estrutura dos dados que o cliente enviará para adicionar ou deletar
type FileInfo struct {
	IP       string `json:"ip"`
	FileName string `json:"filename"`
	Hash     int    `json:"hash"`
	Action   string `json:"action"` // Pode ser "add" ou "delete"
}

// Estrutura para consulta de hash
type HashQuery struct {
	Hash int `json:"hash"`
}

// Mapa global que mapeia o IP para um mapa de arquivos (nome do arquivo -> hash)
var fileMap = make(map[string]map[string]int)
var mutex sync.Mutex

func main() {
	// Escuta na porta 8000
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Servidor ouvindo na porta 8000...")

	for {
		// Aceita uma conexão de um cliente
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // falhas na conexão
			continue
		}
		// Trata a conexão em uma goroutine para permitir múltiplas conexões
		go handleConn(conn)
	}
}

// handleConn trata a conexão recebida de um cliente
func handleConn(c net.Conn) {
	defer c.Close()

	// Buffer para armazenar os dados recebidos
	var buf = make([]byte, 4096) // Aumenta o buffer para suportar mensagens maiores

	// Lê os dados enviados pelo cliente
	n, err := c.Read(buf)
	if err != nil {
		if err != io.EOF {
			log.Println("Erro ao ler dados:", err)
		}
		return
	}

	// Log para depuração
	log.Printf("Dados recebidos: %s\n", string(buf[:n]))

	// Deserializa o JSON enviado pelo cliente
	var fileInfo FileInfo
	err = json.Unmarshal(buf[:n], &fileInfo)
	if err == nil && (fileInfo.Action == "add" || fileInfo.Action == "delete") {
		// Atualiza a estrutura de dados do servidor com o IP, nome do arquivo e hash
		updateFileMap(fileInfo)
		_, err = c.Write([]byte("Arquivos e hashes atualizados com sucesso\n"))
		if err != nil {
			log.Println("Erro ao enviar resposta:", err)
		}
		return
	}

	// Se não for um FileInfo, tenta interpretar como uma consulta de hash
	var query HashQuery
	err = json.Unmarshal(buf[:n], &query)
	if err == nil && query.Hash != 0 {
		// Consulta a lista de IPs que possuem o hash
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

// updateFileMap atualiza o mapa global de arquivos e hashes
func updateFileMap(fileInfo FileInfo) {
	mutex.Lock()
	defer mutex.Unlock()

	// Verifica se o IP já existe no mapa, senão, inicializa o mapa de arquivos
	if _, exists := fileMap[fileInfo.IP]; !exists {
		fileMap[fileInfo.IP] = make(map[string]int)
	}

	if fileInfo.Action == "add" {
		// Adiciona o arquivo e hash para o IP
		fileMap[fileInfo.IP][fileInfo.FileName] = fileInfo.Hash
	} else if fileInfo.Action == "delete" {
		// Remove o arquivo para o IP
		delete(fileMap[fileInfo.IP], fileInfo.FileName)
		// Opcional: Remover o IP se não houver mais arquivos associados a ele
		if len(fileMap[fileInfo.IP]) == 0 {
			delete(fileMap, fileInfo.IP)
		}
	}

	log.Printf("Mapeamento atualizado para IP %s: %+v\n", fileInfo.IP, fileMap[fileInfo.IP])
}

// getIPsForHash retorna uma lista de IPs que possuem o hash fornecido
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
