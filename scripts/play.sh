#!/bin/bash

curl localhost:8000
curl localhost:8000/health

curl localhost:8000/api/v1/users

curl -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"nam", "password":"nam"}'
