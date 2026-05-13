# Perguntas

### Quais são os principais desafios de controle de concorrência no Raft?
Olhando em uma perspectiva geral o maior desafio de controle de concorrência é a própria latência da comunicação pois é o calcanhar de Aquiles em qualquer sistema distribuído, principalmente os que dependem de timeouts. Porém se formos olhar em uma perspectiva interna o maior desafio de concorrência é sincronização de estados e sua atomicidade. Sem boas definições em como lidar com `AppendEntries` de um líder que acaba de ser deposto os estados podem acabar dessincronizados.  

### Como é o tipo de consistência no Raft?
A consistência no Raft é classificada como `strong consistency` pois abre mão da disponibilidade (AP) para que as operações sejam sobre as informações mais atualizadas o possível (CP).

### Como o algoritmo evita inconsistências nos dados mesmo com falhas e mensagens concorrentes?
A inconsistência é evitada pois:
- Líder força `followers` a entrarem em consenso
- Líder só é eleito se estiver atualizado
- Entradas se tornam `committed` apenas no termo atual

### Quais problemas podem ocorrer se temporizadores e threads não forem corretamente sincronizados na implementação?
De forma mais geral os problemas podem aparecer como:
- Corrida de `timer` e RPCs gerando `timeout` falso
- Liderança simultânea
- Livelock (ex.: troca constante de líder, sem progresso)
- Deadlocks
