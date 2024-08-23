#!/bin/bash

curl -X POST http://localhost:8080/migrate-webhook \
  -H "Content-Type: application/json" \
  -d '{
    "repository": {
      "name": "repository-name",
      "clone_url": "https://github.com/username/repository-name.git"
    }
  }'
