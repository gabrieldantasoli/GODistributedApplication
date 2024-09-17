package main

import (
	"encoding/json"
	"fmt"
	"hybridp2p/networks"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Assinatura struct {
	Sum int    `json:"sum"`
	IP  string `json:"ip"`
}

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
func sum(filePath string, out chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	data, err := readFile(filePath)
	if err != nil {
		log.Printf("Erro ao ler o arquivo %s: %v", filePath, err)
		return
	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}

	//Identifica o ip da máquina
	ip, err := networks.GetLocalIP()
	if err != nil {
		log.Fatalf("Erro ao obter IP: %v", err)
	}

	assinatura := Assinatura{Sum: _sum, IP: ip}

	jsonData, err := json.Marshal(assinatura)
	if err != nil {
		log.Printf("Erro ao serializar JSON: %v", err)
		return
	}

	out <- string(jsonData)
}

// print the totalSum for all files and the files with equal sum
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	var totalSum int64
	sums := make(map[int][]string)
	sumsChannel := make(chan string)
	var wg sync.WaitGroup

	for _, path := range os.Args[1:] {
		wg.Add(1)
		go sum(path, sumsChannel, &wg)
	}

	go func() {
		wg.Wait()
		close(sumsChannel)
	}()

	for jsonStr := range sumsChannel {
		var v Assinatura
		if err := json.Unmarshal([]byte(jsonStr), &v); err != nil {
			log.Printf("Erro ao desserializar JSON: %v", err)
			continue
		}

		sums[v.Sum] = append(sums[v.Sum], v.IP)
		totalSum += int64(v.Sum)
	}

	fmt.Println("Total Sum: ", totalSum)

	for sum, files := range sums {
		fmt.Printf("Sum %d: %v\n", sum, files)
	}
}
