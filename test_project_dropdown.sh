#!/bin/bash

echo "🧪 Testing Project Dropdown API"
echo "================================"
echo ""

# Start the server in background
echo "Starting server..."
go run cmd/agent-house/main.go --serve --port 8081 &
SERVER_PID=$!

# Wait for server to start
sleep 3

echo ""
echo "📡 Testing /api/projects endpoint..."
echo ""

# Test the API
response=$(curl -s http://localhost:8081/api/projects)

if [ $? -eq 0 ]; then
    echo "✅ API Response:"
    echo "$response" | python3 -m json.tool 2>/dev/null || echo "$response"

    echo ""
    echo "📊 Project Summary:"
    echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4 | while read project; do
        task_count=$(echo "$response" | grep -o "\"id\":\"$project\"" -A 5 | grep -o '"task_count":[0-9]*' | cut -d':' -f2)
        echo "  - $project ($task_count tasks)"
    done
else
    echo "❌ Failed to connect to API"
fi

echo ""
echo "Stopping server..."
kill $SERVER_PID 2>/dev/null

echo "✅ Test complete!"
