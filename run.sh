#!/bin/bash

# WebGL Water Tutorial Go Port - Run Script
set -e

echo "🌊 WebGL Water Tutorial Go Port"
echo "==============================="

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: Please run this script from the go-port directory"
    exit 1
fi

# Kill any existing server processes
pkill -f "webgl-water-server" 2>/dev/null || true
pkill -f "server.*8080" 2>/dev/null || true

# Build the application
echo "🔨 Building application..."
go build -o server ./cmd/server

# Set environment variables
export PORT=${PORT:-8080}
export ASSETS_PATH="./assets"
export STATIC_PATH="./web/static"

echo "🚀 Starting server on port $PORT..."
echo "📱 Open: http://localhost:$PORT"
echo "⏹️  Press Ctrl+C to stop"
echo ""

# Run the server
./server -port "$PORT" -assets "$ASSETS_PATH" -static "$STATIC_PATH"
