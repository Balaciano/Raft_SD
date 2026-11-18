# Raft — Parte 2A (Sistemas Distribuídos)

Implementação em Go da eleição de líder e dos heartbeats do Raft, conforme a parte 2A do laboratório.

**Integrantes:** Eduardo Osborne, João Pedro Balaciano e Rafaela Fiuza

## Estrutura

| Caminho | Conteúdo |
| --- | --- |
| `src/raft/raft.go` | Nossa implementação da parte 2A |
| `src/raft/election_extra_test.go` | 11 testes adicionais escritos pelo grupo |
| `src/raft/test_test.go`, `config.go`, `persister.go`, `util.go` | Arquivos fornecidos, sem alteração |
| `src/labrpc/` | Biblioteca de RPC fornecida, sem alteração |
| `scripts/testar.sh` | Executa os testes 2A no macOS/Linux e salva os logs em `resultados/` |
| `scripts/testar.ps1` | Mesmo script para Windows (PowerShell) |
| `scripts/verificar_base.py` | Confere se os arquivos fornecidos estão idênticos aos originais |
| `scripts/empacotar.py` | Gera um `.zip` com o código, os scripts e a documentação |
| `docs/RELATORIO.md` | Relatório do trabalho |

## Como testar

Requisito: Go instalado (`go version` deve funcionar). Na raiz do projeto:

```bash
bash scripts/testar.sh            # testes 2A uma vez
bash scripts/testar.sh 20         # 20 repetições
bash scripts/testar.sh 5 race     # detector de corridas, 5 repetições
```

O script define `GOPATH` na raiz do projeto e `GO111MODULE=off`, necessários porque o código-base usa o modo GOPATH. Para rodar sem o script:

```bash
export GOPATH="$PWD" GO111MODULE=off
go test -v -run 2A raft
go test -race -run 2A raft
```

No Windows, use `.\scripts\testar.ps1`, `.\scripts\testar.ps1 -Repeticoes 20` e `.\scripts\testar.ps1 -Race -Repeticoes 5`.

## Resultados

Em um MacBook Air (Apple M1, Go 1.27.4), todas as 26 rodadas da suíte 2A passaram, incluindo 5 com `-race`, sem corridas de dados. A reeleição levou no máximo 1,442 s (limite de 5 s) e a taxa de heartbeats ficou entre 7,5 e 8,0 por segundo (limite de 10). Detalhes em `docs/RELATORIO.md`.

## Escopo

Somente a parte 2A. `Start` mantém a assinatura do esqueleto e retorna índice `-1`, pois a replicação de comandos (2B) e a persistência (2C) não fazem parte desta entrega. Por isso, `go test` sem o filtro `-run 2A` falha nos testes 2B e 2C.
