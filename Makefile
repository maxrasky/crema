.PHONY: test
test:
	@docker compose up -d memcached
	@go test ./... -cover -race -count=1

.PHONY: fmt
fmt:
	@gofmt -l -w -s .
	@goimports -w .

DIR=internal/service
.PHONY: proto
proto:
	protoc -I="$(DIR)" --go_opt=paths=source_relative --go_out="$(DIR)" "$(DIR)/proto/service.proto"
	protoc -I="$(DIR)" --go-grpc_opt=require_unimplemented_servers=false --go-grpc_opt=paths=source_relative --go-grpc_out="$(DIR)" "$(DIR)/proto/service.proto"

.PHONY: run
run:
	docker compose up

.PHONY: drop
drop:
	docker rmi -f crema-app


