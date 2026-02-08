#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

echo "Pulling latest image..."
docker compose pull

echo "Restarting with new image..."
docker compose up -d

echo "Done."
docker compose ps
