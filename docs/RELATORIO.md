# Trabalho Prático 2 — Raft (Parte 2A)

**Disciplina:** Sistemas Distribuídos — 2026/2

**Integrantes:** Eduardo Osborne, João Pedro Balaciano e Rafaela Fiuza

## 1. Introdução

Neste trabalho, implementamos a parte 2A do Raft em Go: eleição de líder e heartbeats, usando o pacote `labrpc` fornecido. Nossa implementação elege um único líder, mantém esse líder enquanto não há falhas e elege outro quando ele cai ou fica isolado. A replicação de log (2B) e a persistência (2C) ficaram fora do escopo, como pede o enunciado.

## 2. Implementação

O código está em `src/raft/raft.go` e segue a Figura 2 do artigo estendido. Cada servidor guarda o próprio termo, o voto e o papel (seguidor, candidato ou líder), protegidos por um mutex, e se comunica com os outros apenas por RPC.

- **Eleição:** quando o timeout expira, o servidor vira candidato, incrementa o termo e pede votos. Com a maioria, vira líder.
- **RequestVote:** cada servidor vota em no máximo um candidato por termo e recusa candidatos com termo antigo ou log desatualizado.
- **AppendEntries (heartbeat):** o líder envia heartbeats periódicos. Quem recebe volta a ser seguidor e reinicia o timeout.
- **Termo maior:** qualquer servidor que vê um termo maior que o seu volta a ser seguidor.

Usamos heartbeat a cada 120 ms, para ficar abaixo do limite de 10 por segundo, e timeout de eleição sorteado entre 450 e 900 ms. O sorteio evita empates e permite várias rodadas de eleição antes de 5 s.

## 3. Política de colaboração

Desenvolvemos o código a partir do enunciado, do artigo e do código-base da disciplina, sem consultar soluções de outros grupos ou da internet. Não alteramos os arquivos fornecidos. Dividimos o trabalho assim:

| Integrante | Responsabilidades |
| --- | --- |
| João Pedro Balaciano | Estado do Raft, RPC `RequestVote` e lógica de eleição e contagem de votos |
| Rafaela Fiuza | RPC `AppendEntries`, envio de heartbeats pelo líder e temporizadores |
| Eduardo Osborne | Testes, ajustes e relatório |

O estudo do artigo e a revisão final do código foram feitos por todo o nosso grupo.

## 4. Ferramentas utilizadas

Usamos Go 1.27.4 (implementação, testes e detector de corridas), o `labrpc` fornecido e Git. Também usamos o ChatGPT e o Claude para gerar trechos de código, tirar dúvidas, encontrar soluções para problemas e redigir este relatório. Revisamos e testamos tudo o que foi produzido com essas ferramentas.

## 5. Dificuldades encontradas

- **Ambiente:** o código-base é de 2018 e só compila no Go atual com `GOPATH` na raiz do projeto e `GO111MODULE=off`. Resolvemos isso no script `testar.sh`.
- **Tempos:** o limite de 10 heartbeats por segundo impede usar os 150–300 ms do artigo. Ajustamos os tempos até equilibrar esse limite com a reeleição em menos de 5 s.
- **Respostas atrasadas:** um voto pode chegar depois que o candidato já mudou de termo. Passamos a conferir o termo e o papel antes de contar cada voto.

## 6. Lições aprendidas

- Segurança e disponibilidade dependem de mecanismos diferentes: o voto único e a maioria garantem um só líder por termo, e os timeouts só afetam a rapidez da eleição.
- Sistemas distribuídos não são determinísticos. Um teste que passa uma vez pode falhar na próxima, por isso repetimos a suíte várias vezes e usamos o detector de corridas.

## 7. Resultados

Rodamos os testes em 18/11/2026, num MacBook Air (Apple M1, 8 GB, macOS 26.0.1, Go 1.27.4). A suíte 2A tem os 2 testes oficiais e 11 testes adicionais que escrevemos.

| Execução | Rodadas aprovadas | Testes aprovados |
| --- | --- | --- |
| `testar.sh 1` | 1/1 | 13/13 |
| `testar.sh 20` | 20/20 | 260/260 |
| `testar.sh 5 race` | 5/5 | 65/65, sem corridas de dados |
| **Total** | **26/26** | **338/338** |

- **Reeleição (limite de 5 s):** entre 0,351 s e 1,442 s, com média de 0,566 s em 26 medições.
- **Heartbeats (limite de 10 por segundo):** entre 7,5 e 8,0 por segundo.

## 8. Conclusão

Nossa implementação cumpre os requisitos da parte 2A. Ela passou em todas as 26 rodadas, sem corridas de dados, e ficou com folga nos dois limites do enunciado. A estrutura já tem a base para a parte 2B, como as entradas de log e a comparação de logs na votação.

## Referências

1. Enunciado do Trabalho Prático: Raft — Sistemas Distribuídos, 2026/2.
2. ONGARO, D.; OUSTERHOUT, J. *In Search of an Understandable Consensus Algorithm (Extended Version)*. 2014.
3. ONGARO, D.; OUSTERHOUT, J. *In Search of an Understandable Consensus Algorithm*. USENIX ATC, 2014.
