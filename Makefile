MIGRATIONS_PATH = common/db/migrations
CONFIG_PATH = ./config/local.yaml
POSTGRES_URL = postgres://postgres:postgres@localhost:5432/restaurant?sslmode=disable

.PHONY: run-sso
run-auth:
	@go run sso/cmd/main.go --config=$(CONFIG_PATH)

.PHONY: gen-proto
gen-proto:
	@name=$(name);
	@protoc --proto_path=common/api/proto \
	  --go_out=common/api/gen --go_opt=paths=source_relative \
		--go-grpc_out=common/api/gen --go-grpc_opt=paths=source_relative \
		common/api/proto/$(name)/*.proto

.PHONY: migrate-create
migrate-create:
	@name=$(name);
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(name)

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(POSTGRES_URL) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(POSTGRES_URL) down $(filter-out $@,$(MAKECMDGOALS))