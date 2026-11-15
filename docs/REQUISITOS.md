# Interpretação dos requisitos

Os PDFs são fontes de requisitos e fundamentação. O pedido do usuário define
a entrega: interpretar o material, implementar o projeto e deixar resultados
dependentes de execução como **FALTA COMPLETAR**.

O PDF de sete páginas é o enunciado; os demais são artigos de referência.
As páginas 1-4 apresentam o laboratório e sua interface geral. A página 4
exige aprovação em 2A. A seção “Tarefa”, p. 5, pede eleição e heartbeats.
A saída ilustrativa que também mostra 2B não estabelece uma tarefa adicional
de replicação. A entrega, portanto, implementa 2A.

| Requisito | Fonte | Atendimento / evidência pendente |
| --- | --- | --- |
| Go e instâncias por participante | Enunciado, p. 2 | `Raft`, `Make` |
| Comunicação somente por RPC | Enunciado, p. 2 e 5 | `labrpc.ClientEnd.Call`; estado próprio por instância |
| Não usar arquivos no algoritmo | Enunciado, p. 2 | Sem operações de arquivo em `raft.go` |
| Base fornecida e labrpc original | Enunciado, p. 2-3 e 5 | Seis arquivos conferidos por hash Git |
| Interface preservada | Enunciado, p. 4-5 | `Make`, `Start`, `GetState`, `ApplyMsg` |
| Eleição inicial, manutenção e reeleição | Enunciado, p. 5 | Votos, termos, maioria e temporizadores; testes: FALTA COMPLETAR |
| Campos públicos nas RPCs | Enunciado, p. 5 | Nomes iniciados por maiúsculas |
| Timeouts distintos | Enunciado, p. 6 | Sorteio por instância em [450 ms, 900 ms) |
| Até 10 rodadas de heartbeat/s | Enunciado, p. 6 | Intervalo programado mínimo de 120 ms; medição: FALTA COMPLETAR |
| Reeleição em até 5 s com maioria | Enunciado, p. 6 | Rodadas aleatórias; tempo observado: FALTA COMPLETAR |
| Verificar corridas | Enunciado, p. 7 | Scripts com `-race`; resultado: FALTA COMPLETAR |
| Relatório obrigatório | Enunciado, p. 1 | `RELATORIO.md`, seções 4-7 |
| Código e relatório no Classroom | Enunciado, p. 7 | Arquivos locais prontos para revisão; envio pelo grupo |

Os scripts de teste gravam logs fora do algoritmo. Esses arquivos não são
usados como armazenamento ou comunicação pelos participantes Raft.

A Figura 2 (p. 4 do artigo estendido) fundamenta as transições e RPCs. As
seções 5.1-5.2 (p. 5-6) explicam termos, papéis, maioria e timeout aleatório.
A seção 5.4.1 (p. 8) fundamenta a comparação dos logs. A seção 5.6 (p. 10)
explica a relação entre comunicação e temporizadores. As regras completas
de replicação da Figura 2 pertencem às próximas etapas.
