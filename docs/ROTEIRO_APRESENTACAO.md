# Roteiro de estudo e apresentação

## 1. Problema e escopo

Explique que Raft coordena réplicas por meio de um log, mas esta etapa trata
somente de escolher e manter o líder. Abra `src/raft/raft.go` e mostre os três
papéis. Não apresente `Start` como replicação implementada.

## 2. Eleição passo a passo

1. Todos começam seguidores e sorteiam seus próprios timeouts.
2. Um timeout expira: o nó incrementa o termo e vota em si.
3. `RequestVote` é enviado por `labrpc` a cada outro participante.
4. Cada destinatário verifica termo, voto anterior e atualização do log.
5. Maioria dos votos transforma o candidato em líder.
6. O líder envia um heartbeat imediatamente e repete periodicamente.

Exemplo com cinco nós: três votos são necessários. Dois votos não bastam,
mesmo que apenas dois servidores estejam acessíveis.

## 3. Falhas e respostas atrasadas

Ao desconectar o líder, os seguidores deixam de receber heartbeats e iniciam
novas eleições. A maioria conectada pode eleger outro líder. Ao reconectar,
o líder antigo abandona seu papel ao descobrir um termo superior.

Uma resposta de uma eleição anterior não deve ser contada na rodada atual.
Mostre a verificação do termo e do papel nos callbacks de `startElectionLocked`.
Mostre também que um termo superior é observado antes de descartar o voto.

## 4. Perguntas que os integrantes devem saber responder

- Por que a maioria evita dois líderes eleitos no mesmo termo?
- Por que o voto não é apagado ao receber heartbeat no mesmo termo?
- Por que pedidos de voto recusados não reiniciam o timeout?
- Por que usar timeout aleatório em vez de um tempo fixo para todos?
- Por que o mutex é liberado durante a espera de RPC?
- O que ocorre com uma minoria isolada? E com um líder antigo isolado?
- Quais mudanças seriam necessárias para implementar 2B e 2C?
- O que os testes realmente demonstraram e o que não foi medido?

## 5. Demonstração

Use os scripts do README e mostre os logs reais dos testes oficiais, dos testes
adicionais, das medições e do detector de corridas. Preencha previamente
`RESULTADOS.md`. **FALTA COMPLETAR:** inserir aqui os resultados observados
e um exemplo real de dificuldade resolvida pelo grupo.
