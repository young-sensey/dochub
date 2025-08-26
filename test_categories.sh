#!/bin/bash

echo "Testing Category API endpoints..."

# Создание категории
echo "Creating category..."
CREATE_RESPONSE=$(curl -s -X POST http://localhost/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Техническая документация", "description": "Документы технического характера"}')

echo "Create response: $CREATE_RESPONSE"

# Получение списка категорий
echo "Getting categories list..."
LIST_RESPONSE=$(curl -s -X GET http://localhost/categories)
echo "List response: $LIST_RESPONSE"

# Получение конкретной категории (если есть ID)
CATEGORY_ID=$(echo $CREATE_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
if [ ! -z "$CATEGORY_ID" ]; then
    echo "Getting category with ID: $CATEGORY_ID"
    GET_RESPONSE=$(curl -s -X GET http://localhost/categories/$CATEGORY_ID)
    echo "Get response: $GET_RESPONSE"
    
    # Обновление категории
    echo "Updating category..."
    UPDATE_RESPONSE=$(curl -s -X PUT http://localhost/categories/$CATEGORY_ID \
      -H "Content-Type: application/json" \
      -d '{"name": "Обновленная техническая документация", "description": "Обновленное описание"}')
    echo "Update response: $UPDATE_RESPONSE"
fi

echo "Category API test completed." 