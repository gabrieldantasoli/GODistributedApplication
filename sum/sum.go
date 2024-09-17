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
func sum(filePath string, outJson chan []byte) {
	data, _ := readFile(filePath)
	// 	if err != nil {
	// 		return 0, err
	// 	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}
	//Identifica o ip da máquina
	ip, err := networks.GetLocalIP()
	if err != nil {
		log.Fatalf("Erro ao obter IP: %v", err)
	}

	jsonMap := map[string]interface{}{
		"Sum": _sum,
		"Ip":  ip,
	}

	jsonData, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}
	outJson <- jsonData
}

// print the totalSum for all files and the files with equal sum
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	jsonChannel := make(chan []byte)
	var wg sync.WaitGroup

	for _, path := range os.Args[1:] {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			sum(path, jsonChannel)
		}(path)
	}

	go func() {
		wg.Wait()
		close(jsonChannel)
	}()

	for jsonData := range jsonChannel {
		fmt.Println(string(jsonData))
	}

}
