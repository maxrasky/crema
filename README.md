---
Before running application one has to have memcached service running.

It can be achieved by running `docker compose up -d`

---
App has `Makefile`, which can be utilized for running tests and application itself.

- `make test` - runs tests
- `make run` - runs application as gRPC server
