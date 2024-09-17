package main

import (
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

// Estrutura dos dados que o cliente enviará para adicionar ou deletar
type FileInfo struct {
	IP       string `json:"ip"`
	FileName string `json:"filename"`
	Hash     int    `json:"hash"`
	Action   string `json:"action"` // Pode ser "add" ou "delete"
}

func main() {
	// Caminho da pasta a ser monitorada
	watchDir := "./dataset"

	log.Println("Cliente Iniciado...")

	// Inicia o monitoramento da pasta
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Inicia uma goroutine para tratar eventos do watcher
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// Trata diferentes tipos de eventos
				if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
					handleFileEvent(event.Name, "add")
				} else {
					handleFileEvent(event.Name, "delete")
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Erro:", err)
			}
		}
	}()

	// Adiciona a pasta ao watcher
	err = watcher.Add(watchDir)
	if err != nil {
		log.Fatal(err)
	}

	// Mantém o programa rodando
	select {}
}

// handleFileEvent trata eventos de criação, modificação ou remoção de arquivos
func handleFileEvent(filePath string, action string) {
	ip := getLocalIP()

	// Se for um arquivo que está sendo criado ou modificado, calcula o hash
	var fileHash int
	if action == "add" {
		hash, err := sum(filePath)
		if err != nil {
			log.Println("Erro ao calcular hash:", err)
			return
		}
		fileHash = hash
	} else if action == "delete" {
		// Para exclusão, o hash não é necessário, então não calcula
		fileHash = 0
	}

	// Prepara os dados a serem enviados ao servidor
	fileInfo := FileInfo{
		IP:       ip,
		FileName: getFileName(filePath),
		Hash:     fileHash,
		Action:   action,
	}

	// Envia a informação ao servidor
	err := sendToServer(fileInfo)
	if err != nil {
		log.Println("Erro ao enviar dados para o servidor:", err)
	}
}

// sendToServer envia os dados de hash para o servidor
func sendToServer(fileInfo FileInfo) error {
	// Conecta ao servidor
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		return err
	}
	defer conn.Close()

	// Converte a estrutura FileInfo para JSON
	data, err := json.Marshal(fileInfo)
	if err != nil {
		return err
	}

	// Envia os dados ao servidor
	_, err = conn.Write(data)
	return err
}

// getFileName retorna o nome do arquivo a partir do caminho completo
func getFileName(filePath string) string {
	parts := strings.Split(filePath, "/")
	return parts[len(parts)-1]
}

// getLocalIP retorna o endereço IP local da máquina
func getLocalIP() string {
	// Este é um exemplo simples, adapte conforme necessário para obter o IP correto
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return "localhost"
}
