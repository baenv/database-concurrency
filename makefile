#!make
include .env
export $(shell sed 's/=.*//' .env)

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down
	
.PHONY: force-remove-img
force-remove-img: down
	docker rmi -f indexer && docker rmi -f server && docker rmi -f queue

.PHONY: up-latest
up-latest: force-remove-img up

.PHONY: ent-gen
ent-gen: 
	go generate ./ent 
	
.PHONY: migrate-new
migrate-new:
	go run -mod=mod ent/migrate/main.go postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@127.0.0.1:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable ${name}
	
.PHONY: migrate-latest
migrate-latest:
	atlas migrate apply \
  --dir "file://ent/migrate/migrations" \
	-u "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@127.0.0.1:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"
	
.PHONY: migrate-status
migrate-status:
	atlas migrate status \
  --dir "file://ent/migrate/migrations" \
	-u "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@127.0.0.1:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"
	
.PHONY: ent-init
ent-init: 
	go run -mod=mod entgo.io/ent/cmd/ent init ${name}
