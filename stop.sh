#!/bin/bash
# Orchestra - Stop all services

cd "$(dirname "$0")/backend/deployments/production"

echo "🛑 Stopping Orchestra services..."
docker compose down

echo ""
echo "✅ All services stopped!"
