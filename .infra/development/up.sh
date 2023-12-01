export IMAGE=money-keeper-backend:local-$(uuidgen)

BUILD_CONTEXT=./
DOCKERFILE=.infra/Dockerfile_dev

echo "Building image: $IMAGE"

docker build -f $DOCKERFILE $BUILD_CONTEXT -t $IMAGE
docker compose -f .infra/development/docker-compose.yml up -d