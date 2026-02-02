#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

echo "🔄 Pulling latest changes..."
git pull origin main

echo "🚀 Rebuilding (with cache)..."
DOCKER_BUILDKIT=1 docker compose build
docker compose up -d

echo "✅ Done."
docker compose ps
