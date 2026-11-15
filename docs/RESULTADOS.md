# Ficha de execução: FALTA COMPLETAR

Preencher somente a partir dos logs e do ambiente realmente usado.
Os testes Go não foram executados durante a preparação.

## Ambiente

| Campo | Valor |
| --- | --- |
| Responsável / data | FALTA COMPLETAR |
| Sistema operacional | FALTA COMPLETAR |
| Processador / RAM | FALTA COMPLETAR |
| Saída de `go version` | FALTA COMPLETAR |
| Compilador C usado com `-race` | FALTA COMPLETAR |
| Revisão / identificação do código testado | FALTA COMPLETAR |

## Execuções

| Grupo | Comando | Resultado | Duração | Evidência |
| --- | --- | --- | --- | --- |
| Oficiais | `go test -v -run 'Test(InitialElection|ReElection)2A$' -count=1 raft` | FALTA COMPLETAR | FALTA COMPLETAR | FALTA COMPLETAR |
| 2A completo, oficiais e adicionais | `go test -v -run 2A -count=1 raft` | FALTA COMPLETAR | FALTA COMPLETAR | FALTA COMPLETAR |
| Repetição | `go test -run 2A -count=20 raft` | FALTA COMPLETAR: __ / 20 suítes aprovadas | FALTA COMPLETAR | FALTA COMPLETAR |
| Detector de corridas | `go test -race -run 2A -count=5 raft` | FALTA COMPLETAR | FALTA COMPLETAR | FALTA COMPLETAR |

Se o comando agregado com `-count=20` falhar, não deduzir uma contagem de
aprovações sem examinar as execuções. Use `-v` e os logs para contar.
Resultado PASS significa que aquela execução passou; não prova ausência de
todos os defeitos possíveis. Se houver falha, registre-a antes de corrigir.

## Medições já instrumentadas nos testes adicionais

Com GOPATH configurado conforme o README, execute:

```text
go test -v -run 'Test(ReelectionTiming|HeartbeatRate)2A$' -count=10 raft
```

`TestReelectionTiming2A` desconecta o líder em rede confiável com três nós,
mantém dois nós comunicantes e consulta a liderança a cada 10 ms. Registra
`REELEICAO_SEGUNDOS` e falha se não observar novo líder antes de cinco segundos.
É uma medição da detecção pelo teste, com resolução de consulta e influência
do escalonador. Não representa garantia de tempo sob quaisquer perdas de rede.

`TestHeartbeatRate2A` espera estabilização e conta RPCs recebidas por um seguidor
durante aproximadamente dois segundos. Sem comandos ou eleição no intervalo,
essas RPCs correspondem a heartbeats. O teste verifica estabilidade do termo,
registra `HEARTBEATS_RPC`, `JANELA_SEGUNDOS`, `TAXA_RPC_POR_SEGUNDO` e exige taxa
positiva de até 10 por segundo. Essa medida é média por seguidor na janela,
não uma contagem agregada do cluster nem prova de todos os intervalos de envio.

| Amostra | Reeleição (s) | Heartbeats recebidos | Janela (s) | Taxa (RPCs/s) |
| --- | --- | --- | --- | --- |
| 1 | FALTA COMPLETAR | FALTA COMPLETAR | FALTA COMPLETAR | FALTA COMPLETAR |
| 2 | FALTA COMPLETAR | FALTA COMPLETAR | FALTA COMPLETAR | FALTA COMPLETAR |
| 3 | FALTA COMPLETAR | FALTA COMPLETAR | FALTA COMPLETAR | FALTA COMPLETAR |

Acrescente as outras amostras efetivamente realizadas.

- Menor / média / maior tempo de reeleição: **FALTA COMPLETAR** s.
- Maior taxa média de heartbeat por seguidor observada: **FALTA COMPLETAR** RPCs/s.
- Alertas de corrida e correções: **FALTA COMPLETAR**.
- Falhas funcionais, causas e correções: **FALTA COMPLETAR**.

## Texto para completar no relatório

> Em FALTA COMPLETAR, no ambiente FALTA COMPLETAR, foram executadas FALTA
> COMPLETAR repetições da suíte 2A. Foram aprovadas FALTA COMPLETAR execuções
> e falharam FALTA COMPLETAR. O detector de corridas apresentou FALTA COMPLETAR.
> O tempo máximo de reeleição observado foi FALTA COMPLETAR s, nas condições
> FALTA COMPLETAR. A taxa máxima média observada de heartbeats por seguidor
> foi FALTA COMPLETAR RPCs/s. As evidências estão em FALTA COMPLETAR.

## Conferência final

- [ ] Identificação e colaboração preenchidas com dados reais.
- [ ] IA e demais ferramentas declaradas.
- [ ] Testes oficiais e adicionais executados no código final.
- [ ] Repetições e detector `-race` executados.
- [ ] Falhas corrigidas e testes executados novamente após a última alteração.
- [ ] Medições e evidências preenchidas no relatório.
- [ ] Todos os campos `FALTA COMPLETAR` revisados antes da entrega.
- [ ] Código e relatório enviados no Classroom.
