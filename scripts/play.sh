#!/bin/bash

## health check
# curl localhost:8000
# curl localhost:8000/health


## users
# curl localhost:8000/api/v1/users

# curl -X POST http://localhost:8000/api/v1/users \
#   -H "Content-Type: application/json" \
#   -d '{"email":"nam", "password":"nam"}'

# {"data":{"id":"dccf21e3-903b-461c-934f-3388578c7f76","email":"nam","password":"nam","created_at":1779876246,"updated_at":1779876246},"message":"create new user successfully"}


## targets
# curl localhost:8000/api/v1/targets

# curl -X POST http://localhost:8000/api/v1/targets \
#   -H "Content-Type: application/json" \
#   -d '{"name":"backend", "host":"localhost", "port": "5555"}'

# {"data":{"id":"11d4f170-1f86-4f72-a2a9-8c8f81e6c63e","name":"backend","host":"localhost","port":"5555","created_at":1779876275,"updated_at":1779876275},"message":"create new target successfully"}


## pats
# curl localhost:8000/api/v1/user-pats

# curl -X POST http://localhost:8000/api/v1/user-pats \
#   -H "Content-Type: application/json" \
#   -d '{"user_id":"dccf21e3-903b-461c-934f-3388578c7f76", "target_id":"11d4f170-1f86-4f72-a2a9-8c8f81e6c63e", "ttl_in_hour": 8}'

# {"data":{"plaintext":"nam_tcp_b547ef22c0d6b0dbc94dfcab55aba8f24c5a86b24b3185853e8e8094522982d7","user_pat":{"id":"817ce201-b3d9-44b5-962f-a45619cf1803","user_id":"dccf21e3-903b-461c-934f-3388578c7f76","target_id":"11d4f170-1f86-4f72-a2a9-8c8f81e6c63e","hash_token":"faaa61a84d260e9d64afdce1ffa73ba660fbb9dca3300943eb9830ccb46dece3","created_at":1779876303,"expires_at":1780567503,"revoked_at":0}},"message":"create new user pat successfully"}

## proxy client (save plaintext token from create response)
# go run ./cmd/proxy
# go run ./cmd/client forward -local 127.0.0.1:15432 -proxy localhost:8888 -token "nam_tcp_b547ef22c0d6b0dbc94dfcab55aba8f24c5a86b24b3185853e8e8094522982d7"
# psql -h 127.0.0.1 -p 15432 -U admin mydb
# go run ./cmd/client connect -proxy localhost:8888 -token "nam_tcp_b547ef22c0d6b0dbc94dfcab55aba8f24c5a86b24b3185853e8e8094522982d7"
# go run ./cmd/client send -proxy localhost:8888 -token "nam_tcp_b547ef22c0d6b0dbc94dfcab55aba8f24c5a86b24b3185853e8e8094522982d7" -data "hello"
