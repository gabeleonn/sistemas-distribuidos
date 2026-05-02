# Algoritmo de Consenso Raft
Projeto criado para a disciplina de Sistemas Distribuídos do Programa de Pós-Graduação em Computação Aplicada da UTFPR.

## Estrutura
O projeto é dividido em três partes principais: `proto`, `raft-node` e `raft-client`.

### Proto
Arquivo que define a IDL do projeto, consumida tanto pelo `raft-node` quanto pelo `raft-client`.

### Raft Node
Implementação do algoritmo Raft.
Responsável por executar as diretrizes e regras descritas no paper de 2014 de Ongaro e Ousterhout.

Implementado em Go.

### Raft Cluster
Responsável pelo bootstrap do cluster: criação, inicialização e finalização dos nós.
Também mantém informações globais, como quantidade de nós, atribuição de IDs, portas e operações de reset.
Além disso, atua como camada de observabilidade, consumindo o estado dos nós (via streaming) para fins de debug e inspeção do sistema.

Implementado em Go.

### Raft Client
Cliente responsável por interagir com os nós, enviando comandos para a key-value store e lidando com redirecionamentos para o líder.

Implementado em Python.

## Pastas
```
.
├── grpc/                          # Comunicação via gRPC
│   └── proto/                     # Definições Protocol Buffers (IDL)
│       └── node.proto             # Serviços e mensagens do Raft
├── raft/                          # Implementação do algoritmo Raft
│   ├── models/                    # Modelos de dados do sistema
│   ├── client/                    # Cliente gRPC para interação com os nós
│   └── readme.md                  # Documentação do módulo Raft
└── readme.md                      # Documentação principal do projeto
```

## Instalação das Dependências
WIP

## Como operar?
WIP

## Referências
[https://raft.github.io/raft.pdf](In Search of an Understandable Consensus Algorithm, Ongaro and Ousterhout, 2014)
