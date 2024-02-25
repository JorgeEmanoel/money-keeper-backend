.PHONY: up
up:
	./.infra/development/up.sh

.PHONY: up-prod
up-prod:
	./.infra/production/up.sh

.PHONY: down
down:
	./.infra/development/down.sh

.PHONY: down-prod
down-prod:
	./.infra/down-prod/down.sh

.PHONY: clean
clean: down
	@docker images --format="{{.Repository}}:{{.Tag}}" | grep money | xargs docker rmi

.PHONY: database
database:
	docker exec -it money_keeper_database mysql -u root -proot money_keeper

.PHONY: db
db: database

.PHONY: api
api:
	docker exec -it money_keeper_api sh

.PHONY: logs
logs:
	docker logs -f money_keeper_api

.PHONY: migrate
migrate:
	docker exec -i money_keeper_api sh -c "/app/tmp/main migrate"

.PHONY: migrate-prod
migrate-prod:
	docker exec -i money_keeper_api_prd sh -c "/usr/bin/money-keeper-backend migrate"
