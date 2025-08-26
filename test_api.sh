#!/bin/bash

BASE_URL="http://localhost:8080/dock"

echo "=== Testing DocFlow CRUD API ==="

echo -e "\n1. Creating a document..."
CREATE_RESPONSE=$(curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Document",
    "content": "This is a test document content",
    "author": "John Doe"
  }')
echo $CREATE_RESPONSE | jq '.'

# Extract ID from response
DOC_ID=$(echo $CREATE_RESPONSE | jq -r '.id')
echo "Created document ID: $DOC_ID"

echo -e "\n2. Getting all documents..."
curl -s $BASE_URL | jq '.'

echo -e "\n3. Getting specific document..."
curl -s $BASE_URL/$DOC_ID | jq '.'

echo -e "\n4. Updating document..."
curl -s -X PUT $BASE_URL/$DOC_ID \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Test Document",
    "content": "This is updated content",
    "author": "Jane Smith"
  }' | jq '.'

echo -e "\n5. Getting updated document..."
curl -s $BASE_URL/$DOC_ID | jq '.'

echo -e "\n6. Deleting document..."
curl -s -X DELETE $BASE_URL/$DOC_ID
echo "Document deleted"

echo -e "\n7. Verifying deletion (should return 404)..."
curl -s $BASE_URL/$DOC_ID

echo -e "\n=== Test completed ===" 