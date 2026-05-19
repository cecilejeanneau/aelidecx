#!/usr/bin/env bash
set -euo pipefail

API_KEY_DEFAULT="dev-key-change-me-now-but-long-enough-12345"
PORT=3000
BIND_ADDR=127.0.0.1
API_KEYS="${API_KEYS:-$API_KEY_DEFAULT}"

usage() {
  cat <<EOF
Usage: $0 [-p port] [-b bind_addr] [-k api_key]

Options:
  -p PORT        Port to listen on (default: 3000)
  -b BIND_ADDR   Bind address (use 0.0.0.0 to allow network access)
  -k API_KEY     API key (or set API_KEYS env)
  -h             Show this help
EOF
}

while getopts ":p:b:k:h" opt; do
  case ${opt} in
    p) PORT="$OPTARG" ;;
    b) BIND_ADDR="$OPTARG" ;;
    k) API_KEYS="$OPTARG" ;;
    h) usage; exit 0 ;;
    \?) echo "Invalid option: -$OPTARG" >&2; usage; exit 1 ;;
  esac
done

export APP_ENV=development
export API_KEYS="$API_KEYS"
export PORT="$PORT"
export BIND_ADDR="$BIND_ADDR"

echo "Starting backend: bind=$BIND_ADDR port=$PORT"
cd "$(dirname "$0")/backend"
exec go run -tags sqlite_fts5 .
