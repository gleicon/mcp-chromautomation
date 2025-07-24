#!/bin/bash

# Script to start Chrome with remote debugging enabled
# This allows the MCP service to connect to your existing Chrome instance

echo "Starting Chrome with remote debugging enabled..."

# Detect OS and Chrome location
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    CHROME_PATH="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
    if [ ! -f "$CHROME_PATH" ]; then
        echo "Chrome not found at $CHROME_PATH"
        echo "Install Chrome or update the path in this script"
        exit 1
    fi
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    if command -v google-chrome >/dev/null; then
        CHROME_PATH="google-chrome"
    elif command -v chromium-browser >/dev/null; then
        CHROME_PATH="chromium-browser"
    elif command -v chromium >/dev/null; then
        CHROME_PATH="chromium"
    else
        echo "Chrome/Chromium not found"
        echo "Install with: sudo apt install google-chrome-stable"
        exit 1
    fi
else
    echo "Unsupported OS: $OSTYPE"
    exit 1
fi

# Check if Chrome is already running with debugging
if curl -s http://localhost:9222/json >/dev/null 2>&1; then
    echo "Chrome is already running with debugging enabled on port 9222"
    echo "Your MCP service can now connect to the existing Chrome instance"
    exit 0
fi

# Kill any existing Chrome processes to ensure clean start
echo "Stopping any existing Chrome processes..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    pkill -f "Google Chrome" 2>/dev/null || true
else
    pkill -f chrome 2>/dev/null || true
    pkill -f chromium 2>/dev/null || true
fi

sleep 2

# Start Chrome with debugging enabled
echo "Starting Chrome with remote debugging..."

# Chrome debugging flags
DEBUG_FLAGS=(
    --remote-debugging-port=9222
    --remote-debugging-address=127.0.0.1
    --disable-web-security
    --disable-features=VizDisplayCompositor
    --user-data-dir="$HOME/.chrome-mcp-debug"
)

# Start Chrome in background
if [[ "$OSTYPE" == "darwin"* ]]; then
    "$CHROME_PATH" "${DEBUG_FLAGS[@]}" > /dev/null 2>&1 &
else
    $CHROME_PATH "${DEBUG_FLAGS[@]}" > /dev/null 2>&1 &
fi

# Wait for Chrome to start
echo "Waiting for Chrome to start..."
for i in {1..10}; do
    sleep 1
    if curl -s http://localhost:9222/json >/dev/null 2>&1; then
        echo "Chrome started successfully with debugging enabled!"
        echo ""
        echo "Chrome is now ready for MCP automation"
        echo "Debugging port: http://localhost:9222"
        echo "You can now run: ./build/mcp-chromautomation server"
        echo ""
        echo "Pro tip: Navigate to your favorite sites and log in"
        echo "   Your sessions and cookies will be preserved!"
        exit 0
    fi
    echo "   Attempt $i/10..."
done

echo "Failed to start Chrome with debugging"
echo "Try running Chrome manually with:"
echo "   $CHROME_PATH --remote-debugging-port=9222"
exit 1
