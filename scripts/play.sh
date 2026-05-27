#!/bin/bash

curl localhost:8000
curl localhost:8000/health


## users
curl localhost:8000/api/v1/users

curl -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"nam", "password":"nam"}'


## targets
curl localhost:8000/api/v1/targets

curl -X POST http://localhost:8000/api/v1/targets \
  -H "Content-Type: application/json" \
  -d '{"name":"backend", "host":"localhost", "port": "5000"}'


## pats
curl localhost:8000/api/v1/user-pats

curl -X POST http://localhost:8000/api/v1/user-pats \
  -H "Content-Type: application/json" \
  -d '{"user_id":"f7eacacf-8ecd-4477-ac63-b9942b682686", "target_id":"9edbf175-8b25-4f2b-b721-bb39729e98ba", "ttl_in_hour": 8}'
