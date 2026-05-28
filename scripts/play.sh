#!/bin/bash

## health check
# curl localhost:8000
# curl localhost:8000/health


## users
# curl localhost:8000/api/v1/users

# curl -X POST http://localhost:8000/api/v1/users \
#   -H "Content-Type: application/json" \
#   -d '{"email":"nam", "password":"nam"}'

# {"data":{"id":"d41a6cb7-6633-4c7c-a893-d172d673992c","email":"nam","password":"nam","created_at":1779937649,"updated_at":1779937649},"message":"create new user successfully"}


## targets
# curl localhost:8000/api/v1/targets

# curl -X POST http://localhost:8000/api/v1/targets \
#   -H "Content-Type: application/json" \
#   -d '{"name":"sql", "host":"localhost", "port": "5432"}'

# {"data":{"id":"94570e26-7aaf-467c-810b-c0a495146079","name":"sql","host":"localhost","port":"5432","created_at":1779937675,"updated_at":1779937675},"message":"create new target successfully"}


## pats
# curl localhost:8000/api/v1/user-pats

# curl -X POST http://localhost:8000/api/v1/user-pats \
#   -H "Content-Type: application/json" \
#   -d '{"user_id":"d41a6cb7-6633-4c7c-a893-d172d673992c", "target_id":"94570e26-7aaf-467c-810b-c0a495146079", "ttl_in_hour": 12}'

# {"data":{"plaintext":"nam_tcp_21ebe8e3e4debe1899910017824d8be0b539c0d4323bff830fd3f7b3aa5c84b1","user_pat":{"id":"cf0881ab-0570-42f1-baf2-7f0b9a380510","user_id":"d41a6cb7-6633-4c7c-a893-d172d673992c","target_id":"94570e26-7aaf-467c-810b-c0a495146079","hash_token":"0155dea0bc75e14c8780f5053740fb77480355ab1ff76a612df680e929d59bfb","created_at":1779937785,"expires_at":1780974585,"revoked_at":0}},"message":"create new user pat successfully"}

## proxy client (save plaintext token from create response)
# go run ./cmd/proxy
# go run ./cmd/client forward -local 0.0.0.0:15432 -proxy localhost:8888 -token "nam_tcp_21ebe8e3e4debe1899910017824d8be0b539c0d4323bff830fd3f7b3aa5c84b1"
# psql -h 127.0.0.1 -p 15432 -U admin mydb
# go run ./cmd/client connect -proxy localhost:8888 -token "nam_tcp_21ebe8e3e4debe1899910017824d8be0b539c0d4323bff830fd3f7b3aa5c84b1"
# go run ./cmd/client send -proxy localhost:8888 -token "nam_tcp_21ebe8e3e4debe1899910017824d8be0b539c0d4323bff830fd3f7b3aa5c84b1" -data "hello"
