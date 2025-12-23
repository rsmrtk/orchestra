!/bin/bash
# Orchestra - View logs

cd "$(dirname "$0")/backend/deployments/production"

echo "📝 Viewing logs (Press Ctrl+C to exit)..."
echo ""
docker compose logs -f
