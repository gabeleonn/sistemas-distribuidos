# Checklist — Trabalho Raft

## 1. Requisitos gerais

- [ ] Implementar 4 processos que executam o protocolo Raft.
- [ ] Os processos devem representar os papéis: Seguidor, Candidato, Líder
- [ ] A comunicação entre os processos deve ser feita exclusivamente via gRPC.
- [ ] O cliente deve ser implementado em uma linguagem diferente da utilizada nos processos Raft.
- [ ] Todos os serviços e mensagens devem ser definidos via Protocol Buffers (`.proto`).

---

## 2. Inicialização

- [ ] Inicializar os 4 processos Raft como seguidores.
- [ ] Inicializar um cliente responsável por enviar comandos ao líder.
- [ ] O cliente deve ser capaz de localizar o líder.

---

## 3. Eleição — Valor 8

- [ ] Cada nó inicia como seguidor com um temporizador de eleição aleatório.
- [ ] Ao expirar o temporizador, o nó se torna candidato e inicia uma eleição.
- [ ] Cada nó pode votar apenas uma vez por termo.
- [ ] O líder deve possuir o log atualizado, considerando: Termo e Índice.
- [ ] Um nó não pode votar em um processo que esteja desatualizado.
- [ ] O candidato que obtiver maioria dos votos se torna líder.
- [ ] O líder envia mensagens periódicas de heartbeat usando `AppendEntries`.
- [ ] A ausência de heartbeat leva à detecção de falha e ao início de uma nova eleição.

---

## 4. Replicação — Valor 10

- [ ] O cliente envia dados ao líder via gRPC.
- [ ] O líder adiciona os dados ao seu log.
- [ ] O líder replica as entradas para os seguidores usando `AppendEntries`.
- [ ] Uma entrada é considerada efetivada (`committed`) quando confirmada por maioria dos nós.
- [ ] O líder efetiva a entrada.
- [ ] O líder envia ordem de commit aos seguidores.
- [ ] Na sequência, o líder avisa o cliente.

---

## 5. Interoperabilidade — Valor 3

- [ ] O cliente deve ser implementado em uma linguagem diferente dos nós Raft.
- [ ] Demonstrar que o cliente consegue se comunicar corretamente com o líder via gRPC.
- [ ] Explicar brevemente no relatório como o uso de Protocol Buffers e gRPC possibilita essa interoperabilidade.
