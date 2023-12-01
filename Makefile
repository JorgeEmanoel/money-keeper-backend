.PHONY: up
up:
	./.infra/development/up.sh

.PHONY: down
down:
	./.infra/development/down.sh

.PHONY: clean
clean: down
	@docker images --format="{{.Repository}}:{{.Tag}}" | grep money | xargs docker rmi

.PHONY: database
database:
	docker exec -it money_keeper_database mysql -u root -proot money_keeper

.PHONY: api
api:
	docker exec -it money_keeper_api bash