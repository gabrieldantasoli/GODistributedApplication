package main

import (
	"fmt"
	"hybridp2p/networks"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Assinatura struct {
	sum int
	ip  string
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
func sum(filePath string, out chan<- Assinatura, wg *sysnc.WaitGroup) {
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

	out <- Assinatura{_sum, ip}
}

// print the totalSum for all files and the files with equal sum
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	var totalSum int64
	sums := make(map[int][]string)
	sumsChannel := make(chan Assinatura)
	var wg sync.WaitGroup

	for _, path := range os.Args[1:] {
		wg.Add(1)
		go sum(path, sumsChannel, &wg)
	}

	go func() {
		wg.Wait()
		close(sumsChannel)
	}()

	for v := range sumsChannel {
		sums[v.sum] = append(sums[v.sum], v.ip)
		totalSum += int64(v.sum)
	}

	fmt.Println("Total Sum:", totalSum)

	for sum, files := range sums {
		fmt.Printf("Sum %d: %v\n", sum, files)
	}
}
