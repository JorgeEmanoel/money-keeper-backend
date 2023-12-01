export IMAGE=money-keeper-backend:local-$(uuidgen)

docker compose -f .infra/development/docker-compose.yml down