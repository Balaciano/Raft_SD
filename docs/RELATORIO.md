# Trabalho Prático 2: Raft

**Disciplina:** Sistemas Distribuídos  
**Semestre:** 2026/2  
**Professor(a):** FALTA COMPLETAR  
**Instituição / turma:** FALTA COMPLETAR  
**Prazo indicado no enunciado:** 19 de novembro  
**Data da versão final:** FALTA COMPLETAR

| Integrante | Matrícula | Responsabilidade efetivamente realizada |
| --- | --- | --- |
| FALTA COMPLETAR | FALTA COMPLETAR | FALTA COMPLETAR |

Adicionar linhas para os demais integrantes, respeitando o limite de quatro.

## 1. Objetivo e delimitação

Este trabalho implementa a parte 2A do laboratório Raft: eleição de líder,
manutenção da liderança por heartbeats e reeleição quando a comunicação com
o líder anterior deixa de ocorrer. O código foi escrito em Go e usa o `labrpc`
da base indicada no enunciado. Cada instância mantém estado próprio; não há
variáveis usadas para comunicação direta entre instâncias nem acesso a
arquivos pelo algoritmo.

A seção “Tarefa” do enunciado (p. 5) determina esse escopo e orienta testar
com `go test -run 2A`. Não estão implementados replicação de comandos (2B),
persistência (2C), snapshots ou mudanças de membros. A interface foi mantida:
`Start` retorna índice `-1`, pois não anexa comandos nesta etapa. Não há envio
de comandos a `applyCh`. Esta entrega não é um serviço completo de máquina
de estados replicada.

## 2. Fundamentação

Raft separa eleição, replicação e segurança. Na eleição, um participante pode
ser seguidor, candidato ou líder. Quando seu timeout expira, um seguidor se
torna candidato, incrementa o termo, vota em si e solicita votos. Um candidato
vence ao reunir `floor(N/2) + 1` votos, contando o próprio.

O termo identifica uma época lógica e permite reconhecer mensagens antigas.
Cada participante concede voto a no máximo um candidato por termo. Duas
maiorias do mesmo conjunto necessariamente se intersectam; essa interseção,
com o voto único, impede eleger dois líderes no mesmo termo no modelo sem
reinício de processos considerado em 2A.

O líder transmite `AppendEntries` sem entradas como heartbeats. Um candidato
que recebe heartbeat com termo igual ou maior retorna a seguidor. Pedidos com
termos antigos são recusados. Um termo superior observado em pedido ou resposta
atualiza o termo local, limpa o voto anterior e faz o participante virar seguidor.

Timeouts aleatórios reduzem a chance de empates sucessivos. Sem maioria,
novos candidatos não podem vencer. Um líder isolado pode continuar acreditando
que é líder até conhecer um termo superior. Portanto, a propriedade de segurança
é unicidade de eleição **por termo**, não a ausência instantânea de líderes
autodeclarados em todas as partições.

## 3. Implementação e decisões

### Estado e concorrência

`Raft` guarda termo, voto, papel, votos da rodada e temporizadores. Existe uma
estrutura de entrada de log e uma sentinela no índice zero; em 2A não são
anexados comandos. A regra de atualização dos logs já está presente em
`RequestVote`: comparar primeiro o último termo e depois o último índice.

Um mutex protege o estado mutável e o gerador aleatório da instância. As RPCs
são realizadas em goroutines fora da região crítica, para que um destino lento
não impeça o atendimento de outros pedidos. Após cada resposta, o estado é
reavaliado sob o mutex. Votos só contam enquanto o participante continuar
candidato daquela rodada e a resposta tiver o termo correspondente.

Uma resposta antiga que informar termo superior ainda provoca atualização.
Há uma solicitação de voto por destino a cada rodada, evitando contagem
duplicada. `Kill` é idempotente, encerra o ticker e impede mudanças posteriores
provocadas por callbacks de RPCs. Chamadas já em andamento podem terminar
depois do encerramento, com suas respostas descartadas.

### Regras das RPCs

`RequestVote` recusa termos antigos, candidatos com logs desatualizados e um
segundo candidato no mesmo termo. Somente a concessão do voto reinicia o
timeout; receber um pedido recusado não adia a eleição.

`AppendEntries` recusa termos antigos. Para termo válido, estabelece o estado
de seguidor, reinicia o timeout e verifica o prefixo do log mesmo em heartbeats.
Um heartbeat no mesmo termo preserva o voto já concedido. Pedidos com entradas
não são tratados como replicação bem-sucedida na implementação 2A.

### Temporização

| Parâmetro | Valor | Justificativa |
| --- | --- | --- |
| Heartbeat periódico | mínimo de 120 ms entre rodadas | Aproximadamente 8,33 rodadas/s, abaixo de 10 |
| Timeout de eleição | sorteado em [450 ms, 900 ms) | Espaço entre heartbeats e possibilidade de várias rodadas antes de 5 s |
| Verificação de timers | 10 ms | Verificação periódica sem espera ocupada |

O primeiro heartbeat ocorre ao obter maioria. Os próximos são agendados a
partir do disparo efetivo, evitando rajadas para compensar atrasos. Uma rodada
envia uma mensagem a cada outro participante; a frequência é considerada por
destino/rodada, não como soma de todos os destinos.

Esses valores são decisões do código, não medições. O limite de cinco segundos
para reeleição ainda precisa ser verificado. Perdas, atrasos e escalonamento
podem exigir outras rodadas.

## 4. Política de colaboração

O enunciado exige solução própria e proíbe copiar soluções de outros grupos
ou implementações da rede. Nesta preparação, os arquivos de suporte e testes
foram obtidos exclusivamente da base que o próprio enunciado fornece. Os
seis arquivos preservados foram conferidos por hashes de blob Git e coincidem
com os originais. A implementação de `raft.go` e os testes adicionais foram
produzidos com assistência de IA a partir dos materiais e da API do laboratório.

**Divisão real do trabalho:** FALTA COMPLETAR: quem estudou, implementou,
revisou e executou testes. Não atribuir atividades não realizadas.

**Discussões com outros grupos:** FALTA COMPLETAR: participantes e conceitos
discutidos, ou “não houve”, se aplicável.

**Revisão humana da solução assistida por IA:** FALTA COMPLETAR com a revisão
efetiva e a compreensão dos integrantes sobre o código.

O enunciado orienta não tornar o código público. Os arquivos foram preparados
localmente; a entrega deve seguir o Classroom indicado pela disciplina.

## 5. Ferramentas utilizadas, incluindo IA

| Ferramenta | Uso |
| --- | --- |
| ChatGPT / Codex | Interpretação assistida, implementação, testes adicionais e documentação |
| Conector GitHub | Obtenção e consulta da base oficial; sem consultar soluções prontas |
| Python | Extração dos PDFs e verificação de integridade dos arquivos |
| pypdf / PyMuPDF | Extração textual e conferência visual do enunciado |
| PowerShell | Inspeção do ambiente e preparação dos scripts |
| Git | Inspeção das mudanças locais; sem publicação |
| Go / detector de corridas | FALTA COMPLETAR: versão e uso após configurar o ambiente |

**Declaração de IA:** esta versão do código e do relatório foi produzida com
assistência de ChatGPT/Codex. A equipe deve revisar a lógica, executar os testes
e assumir a responsabilidade pela versão entregue. A revisão e os resultados
finais ainda precisam ser registrados.

## 6. Dificuldades encontradas e tratamento

Uma dificuldade de interpretação foi separar a tarefa 2A da apresentação geral
do laboratório, que também ilustra saídas de 2B. A leitura da seção “Tarefa”
orientou a delimitação explícita adotada no código e na documentação.

Go não foi encontrado no ambiente de preparação, impedindo compilação e testes.
Foram escritos scripts para Windows e Linux/WSL, com repetições e logs.
**FALTA COMPLETAR:** configuração do ambiente, erros observados na execução,
correções realizadas e evidências da execução posterior.

Os desafios técnicos tratados no desenho incluem preservar o voto em mudanças
no mesmo termo, evitar reiniciar o timeout ao recusar votos e descartar respostas
de rodadas encerradas. Essas regras foram incorporadas ao código. Isso não
significa que erros tenham sido observados e corrigidos em uma execução real.

**Dificuldades adicionais do grupo:** FALTA COMPLETAR. Para cada dificuldade,
descrever sintoma, causa, solução e forma de validação.

## 7. Lições aprendidas

A análise evidencia a separação entre segurança e disponibilidade. O voto único
e a maioria sustentam a segurança da eleição; temporizadores e condições de
rede influenciam a rapidez da recuperação.

Outra conclusão de estudo é que respostas de RPC precisam ser interpretadas
com o estado atual, pois termo ou papel podem mudar durante a espera. O mutex
protege memória compartilhada dentro da instância, mas não substitui a validação
lógica da rodada.

A Figura 2 deve ser lida como um conjunto preciso de condições. Reiniciar
timers indiscriminadamente ou apagar votos em qualquer troca de papel pode
contrariar o protocolo.

**Aprendizado efetivo após revisão e execução:** FALTA COMPLETAR com observações
próprias e exemplos. Os parágrafos acima são síntese conceitual para orientar
a reflexão, não um relato fictício de experiências dos integrantes.

## 8. Validação e resultados

São mantidos os dois testes oficiais 2A: eleição inicial e reeleição após falha
de comunicação. Há onze testes adicionais, alguns com múltiplas verificações,
cobrindo voto único, repetição idempotente, RPCs antigas, pedidos recusados,
heartbeat no mesmo termo, termo superior, prefixo do log, encerramento,
servidor único, eleição em rede não confiável e medições de reeleição e heartbeats.

| Verificação | Comando / procedimento | Resultado |
| --- | --- | --- |
| Integridade da base | `python scripts/verificar_base.py` | Seis arquivos oficiais conferidos e idênticos na preparação |
| Testes 2A e compilação | `go test -v -run 2A -count=1 raft` | FALTA COMPLETAR |
| Repetibilidade | `go test -run 2A -count=20 raft` | FALTA COMPLETAR: __ / 20 execuções da suíte aprovadas |
| Corridas de dados | `go test -race -run 2A -count=5 raft` | FALTA COMPLETAR |
| Tempo de reeleição | Medir após desconexão do líder | FALTA COMPLETAR: __ s; limite: 5 s com maioria |
| Frequência de heartbeats | Contar RPCs por seguidor em janela estável | FALTA COMPLETAR: __ RPCs / __ s = __ RPCs/s |

Os comandos pressupõem GOPATH na raiz do projeto e `GO111MODULE=off`.
Os scripts já configuram essas opções e gravam logs em `resultados/`.
Os procedimentos de medição estão em `RESULTADOS.md`.

**Ambiente:** FALTA COMPLETAR: sistema, CPU, RAM, versão Go, compilador C,
data e horário das execuções.

**Situação desta versão:** implementação escrita e integridade da base conferida.
Compilação, aprovação nos testes Go, ausência de corridas e tempos reais
permanecem pendentes. Nenhum resultado de execução Go foi presumido.

## 9. Referências

1. **Trabalho Prático: Raft.** Enunciado fornecido, `Trabalho Raft-Sem2026-2.pdf`,
   7 páginas, semestre 2026/2. Especialmente “Tarefa”, p. 5-7.
2. ONGARO, Diego; OUSTERHOUT, John. **In Search of an Understandable Consensus
   Algorithm (Extended Version).** 20 maio 2014. `raft-extended.pdf`, Figura 2,
   p. 4; seções 5.1, 5.2, 5.4.1 e 5.6.
3. ONGARO, Diego; OUSTERHOUT, John. **In Search of an Understandable Consensus
   Algorithm.** USENIX ATC, 2014. `atc14-paper-ongaro-RaftShort.pdf`.
4. BU-CS451-DS. **goals-2018**, base indicada no enunciado.
   https://github.com/BU-CS451-DS/goals-2018. Referência
   `720120c7f02a9ec1c1221d26c465b4ebc163f014`.
