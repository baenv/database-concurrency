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
