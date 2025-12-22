#!/bin/sh
set -e

echo "Starting Orchestra Frontend..."

# Start nginx in background
echo "Starting nginx..."
nginx

# Start serve in foreground on port 3002
echo "Starting serve on port 3002..."
exec serve -s dist -l 3002
