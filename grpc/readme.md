# Algoritmo de Consenso Raft

Projeto criado para a disciplina de Sistemas Distribuídos do Programa de Pós-Graduação em Computação Aplicada da UTFPR.

O objetivo do projeto é implementar uma versão didática do algoritmo de consenso **Raft**, usando **Go** para os nós do cluster e **gRPC** para a comunicação entre processos. O sistema também possui um cliente em **Python**, responsável por enviar comandos para uma key-value store replicada.

## Visão geral

O projeto é dividido em três partes principais:

- **Raft Node**: implementação de um nó Raft.
- **Raft Cluster**: inicialização e gerenciamento local de múltiplos nós.
- **Raft Client**: cliente Python usado para enviar comandos aos nós.

A comunicação entre os nós acontece via gRPC. Cada nó executa de forma independente, participa da eleição de líder e, após a eleição, pode receber ou replicar comandos conforme as regras do algoritmo Raft.

## Estrutura do projeto

```text
.
├── grpc/
│   ├── application/      # Runtime e process control
│   ├── client/           # Cliente Python gRPC
│   ├── cmd/              # Comandos CLI da aplicação
│   ├── peer/             # Comunicação com outros nós do cluster
│   ├── proto/            # Definições Protocol Buffers
│   ├── raft/             # Implementação do algoritmo Raft
│   ├── service/          # Serviços gRPC expostos pelos nós
│   ├── store/            # Key-value store aplicada sobre o log
│   ├── go.mod
│   ├── go.sum
│   └── main.go
└── readme.md
```

## Requisitos

Para executar o projeto, é necessário ter instalado:

- Go `1.26+`
- Python `3.14+`
- `pip`
- `venv`

## Instalação

### 1. Instalar dependências do Go

A partir da pasta `grpc`:

```bash
go mod tidy
```

### 2. Criar ambiente virtual do cliente Python

A partir da pasta `grpc/client`:

```bash
python3 -m venv .venv
source .venv/bin/activate
```

### 3. Instalar dependências do cliente

Ainda dentro de `grpc/client`:

```bash
python3 -m pip install -r requirements.txt
```

O arquivo `requirements.txt` deve conter:

```text
grpcio
grpcio-tools
```

## Comandos disponíveis

### Cluster

```bash
go run main.go cluster
```

Flags:

```text
-n, --nodes int   Número de nós no cluster Raft
```

Exemplo:

```bash
go run main.go cluster --nodes 4
```

### Node

```bash
go run main.go node
```

Flags:

```text
-i, --id int          ID do nó Raft
-a, --addr string     Endereço do nó Raft
-r, --peers strings   Peers do nó Raft no formato id=addr
```

Exemplo:

```bash
go run main.go node \
  --id 0 \
  --addr localhost:50050 \
  --peers 1=localhost:50051,2=localhost:50052
```

### Client

```bash
python3 client/client.py <host:port> <command>
```

Comandos suportados:

```text
SET:key:value
GET:key
DEL:key
```

Exemplos:

```bash
python3 client/client.py localhost:50052 SET:x:10
python3 client/client.py localhost:50052 GET:x
python3 client/client.py localhost:50052 DEL:x
```

## Referências

- Ongaro, D.; Ousterhout, J. **In Search of an Understandable Consensus Algorithm**. 2014.  
  Disponível em: https://raft.github.io/raft.pdf
