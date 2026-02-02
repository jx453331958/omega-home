#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

echo "🔄 Pulling latest changes..."
git pull origin main

echo "🚀 Rebuilding and restarting..."
docker compose down
docker compose up -d --build

echo "✅ Updated and running."
docker compose ps
