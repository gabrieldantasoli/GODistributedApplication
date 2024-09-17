# GODistributedApplication

# Instruções para Testar a Execução da Aplicação de Registro e Busca de Arquivos

Esta aplicação permite que diferentes máquinas registrem seus endereços IP associados a hashes de arquivos e que clientes busquem máquinas que armazenam um arquivo baseado em seu hash.

## Pré-requisitos

1. **Go instalado**: Certifique-se de que você tem o Go instalado em sua máquina.
   - Para verificar, execute:
     ```bash
     go version
     ```
   - [Download do Go](https://golang.org/dl/)

2. **Compilar o código**:
   - No diretório onde o código está salvo, execute:
     ```bash
     go build -o cliente client.go
     ```

## Execução do Servidor

1. Para iniciar o servidor, execute:
   ```bash
   go run server.go
   ```

## Comandos client
1. Registrar
    ```bash
   ./cliente register <ip> <hash>
   ```
2. Pesquisar
    ```bash
   ./cliente search <hash>
   ```