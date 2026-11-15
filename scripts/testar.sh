#!/usr/bin/env bash
set -euo pipefail
raiz_projeto="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
export GOPATH="$raiz_projeto"
export GO111MODULE=off
repeticoes="${1:-1}"
modo="${2:-normal}"
if ! [[ "$repeticoes" =~ ^[1-9][0-9]*$ ]]; then
  echo 'Uso: bash scripts/testar.sh [repeticoes] [normal|race]' >&2
  exit 2
fi
args=(-v -run 2A "-count=$repeticoes" -timeout=30m)
case "$modo" in
  normal) ;;
  race) args+=(-race) ;;
  *) echo 'Modo deve ser normal ou race.' >&2; exit 2 ;;
esac
mkdir -p "$raiz_projeto/resultados"
arquivo_log="$raiz_projeto/resultados/2A-$modo-$(date +%Y%m%d-%H%M%S).log"
go version | tee "$arquivo_log"
go test "${args[@]}" raft 2>&1 | tee -a "$arquivo_log"
