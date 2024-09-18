# GODistributedApplication
## Explicação Geral

A GODistributedApplication é uma aplicação distribuída que permite o registro e a busca de arquivos entre diferentes máquinas em uma rede. Cada máquina pode registrar o IP associado a hashes de arquivos, permitindo que clientes busquem por esses arquivos através de seus hashes.
Componentes da Aplicação

- **Servidor:** O servidor é responsável por armazenar a associação entre o IP da máquina, o nome do arquivo e seu respectivo hash. A lógica de armazenamento é gerenciada usando um mapa (map), onde o IP da máquina é associado a um outro mapa contendo o nome do arquivo e seu hash.

- **Monitorador de Arquivos:** Este componente monitora automaticamente a pasta dataset no diretório cliente. Ele detecta alterações na pasta, como a adição ou remoção de arquivos. Quando uma alteração ocorre, o monitorador envia uma estrutura de dados para o servidor contendo:
   - O IP da máquina
   - O nome do arquivo
   - O hash do arquivo
   - O tipo de operação (adição ou remoção de arquivo)

Isso permite que o servidor mantenha as informações atualizadas em tempo real, sem que o cliente precise registrar manualmente cada mudança.

- **Cliente:** O cliente permite que o usuário busque por arquivos no servidor usando o hash do arquivo. Ele envia solicitações ao servidor para localizar máquinas que armazenam determinado arquivo.

## Funcionamento Geral

- **Registro de Arquivos:** Quando um arquivo é adicionado à pasta dataset, o monitorador detecta a alteração, gera o hash do arquivo e envia as informações para o servidor.
- **Busca de Arquivos:** O cliente busca por arquivos no servidor utilizando o hash. O servidor retorna o IP da(s) máquina(s) que armazenam o arquivo correspondente.

## Comandos

Antes de executar a aplicação, certifique-se de que o **Go** está instalado na sua máquina:

```bash
go version
```

Build o server e inicie
```bash
go build server.go
```
```bash
./server
```

Em outro terminal, faça o build e inicie o monitoramento
```bash
go build monitoradorDeArquivos.go
```
```bash
./monitoradorDeArquivos
```

Em outro terminal, faça agora o build do client para fazer as buscas
```bash
go build client.go
```
```bash
 ./client search <hash desejado>
```

<br>

- Caso o cliente queira adicionar um arquivo basta ir na pasta dataset e fazer as alteraçoes desejadas.