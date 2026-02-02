#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

ACTION="${1:-up}"

case "$ACTION" in
  up|start)
    echo "🚀 Building and starting Omega Home..."
    docker compose up -d --build
    echo "✅ Done. Visit http://$(hostname -I 2>/dev/null | awk '{print $1}' || echo localhost):${PORT:-3000}"
    ;;
  down|stop)
    echo "🛑 Stopping Omega Home..."
    docker compose down
    echo "✅ Stopped."
    ;;
  restart)
    echo "🔄 Rebuilding and restarting..."
    docker compose down
    docker compose up -d --build
    echo "✅ Restarted."
    ;;
  logs)
    docker compose logs -f --tail=100
    ;;
  status)
    docker compose ps
    ;;
  *)
    echo "Usage: $0 {up|down|restart|logs|status}"
    exit 1
    ;;
esac
