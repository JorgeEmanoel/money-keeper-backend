export IMAGE=money-keeper-backend:production-$(uuidgen)

BUILD_CONTEXT=./
DOCKERFILE=.infra/Dockerfile

echo "Building image: $IMAGE"

docker build -f $DOCKERFILE $BUILD_CONTEXT -t $IMAGE
docker-compose -f .infra/production/docker-compose.yml up -d

echo "Application running. See logs by runnig: make logs-prod"
