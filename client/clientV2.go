package main

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"net"
	"os"
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
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
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
		hash, err := calculateFileHash(filePath)
		if err != nil {
			log.Println("Erro ao calcular hash:", err)
			return
		}
		fileHash = hash
	} else {
		// Para exclusão, não precisamos do hash (a menos que já saibamos ele)
		fileHash = 0 // Pode armazenar em cache ou não, depende da lógica.
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

// calculateFileHash calcula a soma dos bytes de um arquivo
func calculateFileHash(filePath string) (int, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	sum := 0
	for _, b := range data {
		sum += int(b)
	}

	return sum, nil
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
