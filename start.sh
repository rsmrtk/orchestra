#!/bin/bash
# Orchestra - Start all services

cd "$(dirname "$0")/backend/deployments/production"

echo "🚀 Starting Orchestra services..."
docker compose up -d

echo ""
echo "⏳ Waiting for services to be ready..."
sleep 10

echo ""
echo "📊 Service Status:"
docker compose ps

echo ""
echo "✅ Services are running!"
echo ""
echo "🌐 Access the application:"
echo "   Frontend: http://localhost:5252"
echo "   Backend:  http://localhost:8383"
echo "   Database: localhost:2828"
echo ""
echo "📝 View logs:"
echo "   docker compose -f backend/deployments/production/docker-compose.yml logs -f"
echo ""
echo "🛑 Stop services:"
echo "   docker compose -f backend/deployments/production/docker-compose.yml down"
