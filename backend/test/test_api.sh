#!/bin/bash
set -e

ROOT_URL="http://localhost:8000"
API_URL="$ROOT_URL/api"
echo "Testing API at $ROOT_URL"

# 1. Health Check
echo ""
echo "--- 1. Health Check ---"
curl -s "$ROOT_URL/health"
echo ""

# 2. Readiness Check
echo ""
echo "--- 2. Readiness Check ---"
curl -s "$ROOT_URL/health/ready"
echo ""

# 3. CORS Check
echo ""
echo "--- 3. CORS Check ---"
CORS_STATUS=$(curl -s -I -X OPTIONS "$API_URL/items" \
  -H "Origin: http://localhost:8080" \
  -H "Access-Control-Request-Method: GET" | grep -i "Access-Control-Allow-Origin")
echo "$CORS_STATUS"
if [[ -z "$CORS_STATUS" ]]; then
  echo "Error: CORS headers missing"
  exit 1
else
  echo "CORS Headers Verified"
fi

# 4. Create Item
echo ""
echo "--- 4. Create Item ---"
RESPONSE=$(curl -s -X POST "$API_URL/items" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Item", "price": 99.99}')
echo "$RESPONSE" | jq .
ITEM_ID=$(echo "$RESPONSE" | jq -r .id)
echo "Created Item ID: $ITEM_ID"

if [ "$ITEM_ID" == "null" ] || [ -z "$ITEM_ID" ]; then
  echo "Failed to create item"
  exit 1
fi

# 5. Get All Items
echo ""
echo "--- 5. Get All Items ---"
curl -s "$API_URL/items" | jq .

# 6. Get Item by ID
echo ""
echo "--- 6. Get Item by ID ---"
curl -s "$API_URL/items/$ITEM_ID" | jq .

# 7. Update Item
echo ""
echo "--- 7. Update Item ---"
RESPONSE=$(curl -s -X PUT "$API_URL/items/$ITEM_ID" \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Item Name", "price": 150.50}')
echo "$RESPONSE" | jq .

# 8. Delete Item
echo ""
echo "--- 8. Delete Item ---"
curl -s -X DELETE "$API_URL/items/$ITEM_ID" -v 2>&1 | grep "< HTTP/1.1"

# 9. Verify Deletion
echo ""
echo "--- 9. Verify Deletion (Should be 404) ---"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/items/$ITEM_ID")
echo "Status Code: $STATUS"
if [ "$STATUS" -ne 404 ]; then
  echo "Error: Expected 404, got $STATUS"
  exit 1
fi

echo ""
echo "Test Completed Successfully!"
