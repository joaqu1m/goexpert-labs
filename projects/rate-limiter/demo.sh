#!/bin/bash

BASE_URL="http://localhost:8082"

echo "Making 12 requests from same IP to test rate limiting:"
success=0
blocked=0

for i in {1..12}; do
    status=$(curl -s -o /dev/null -w "%{http_code}" -H "X-Real-IP: 10.0.0.100" "$BASE_URL/")
    if [ "$status" -eq 200 ]; then
        ((success++))
        echo -n "✅"
    else
        ((blocked++))
        echo -n "🚫"
    fi
done

echo ""
echo "Result: $success allowed, $blocked blocked"

if [ "$success" -le 10 ] && [ "$blocked" -gt 0 ]; then
    echo "Rate limiting working correctly!"
else
    echo "Rate limiting may not be working"
fi

echo ""
echo "Testing with token (higher limit):"
response=$(curl -s -H "API_KEY: test-token" -H "X-Real-IP: 10.0.0.200" "$BASE_URL/")
echo "Token response: $(echo "$response" | jq -r '.message')"

headers=$(curl -s -I -H "API_KEY: test-token" -H "X-Real-IP: 10.0.0.200" "$BASE_URL/" | grep "X-Ratelimit-Limit")
echo "Token limit: $headers"
