# Validação realizada na preparação

Data de referência: 15/11/2026.

| Verificação | Situação |
| --- | --- |
| Leitura textual dos três PDFs e conferência visual das sete páginas do enunciado | Realizada |
| Comparação de seis arquivos oficiais por SHA-1 de blob Git | Aprovada |
| Sintaxe do script PowerShell pelo parser do PowerShell | Aprovada |
| Sintaxe do script Bash com `bash -n` | Aprovada |
| Leitura UTF-8 dos fontes e documentos | Aprovada |
| Conferência de whitespace do README por `git diff --check` | Aprovada |
| Compilação Go | FALTA COMPLETAR: Go indisponível no ambiente |
| Execução dos testes oficiais e adicionais | FALTA COMPLETAR |
| Detector de corridas `-race` | FALTA COMPLETAR |
| Medições de tempo e frequência | FALTA COMPLETAR |

As verificações estáticas e de integridade não demonstram correção funcional
do Raft. Os resultados pendentes devem ser preenchidos com as execuções reais.
