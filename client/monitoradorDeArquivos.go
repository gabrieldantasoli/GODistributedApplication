package main

import (
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"log"
	"net"
	"os"
	"strings"
)

type FileInfo struct {
	IP       string `json:"ip"`
	FileName string `json:"filename"`
	Hash     int    `json:"hash"`
	Action   string `json:"action"`
}

func main() {
	watchDir := "./dataset"

	log.Println("Cliente Iniciado...")

	// Enviar todos os arquivos da pasta dataset ao iniciar
	sendInitialFiles(watchDir)

	// Iniciar o watcher para monitorar mudanças na pasta
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
					FileEvent(event.Name, "add")
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
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

	err = watcher.Add(watchDir)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}

// Função para enviar os arquivos existentes ao iniciar
func sendInitialFiles(dir string) {
	// Usar os.ReadDir ao invés de ioutil.ReadDir
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Erro ao ler o diretório %s: %v", dir, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			handleFileEvent(dir+"/"+file.Name(), "add")
		}
	}
}

func handleFileEvent(filePath string, action string) {
	ip := getLocalIP()

	var fileHash int
	if action == "add" {
		hash, err := calculateFileHash(filePath)
		if err != nil {
			log.Println("Erro ao calcular hash:", err)
			return
		}
		fileHash = hash
	} else if action == "delete" {
		fileHash = 0
	}

	fileInfo := FileInfo{
		IP:       ip,
		FileName: getFileName(filePath),
		Hash:     fileHash,
		Action:   action,
	}

	err := sendToServer(fileInfo)
	if err != nil {
		log.Println("Erro ao enviar dados para o servidor:", err)
	}
}

func calculateFileHash(filePath string) (int, error) {
	// Usar os.ReadFile ao invés de ioutil.ReadFile
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	sum := 0
	for _, b := range data {
		sum += int(b)
	}

	return sum, nil
}

func sendToServer(fileInfo FileInfo) error {
	conn, err := net.Dial("tcp", "localhost:8000") // Mude para o IP do servidor se necessário
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(fileInfo)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

func getFileName(filePath string) string {
	parts := strings.Split(filePath, "/")
	return parts[len(parts)-1]
}

func getLocalIP() string {
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
