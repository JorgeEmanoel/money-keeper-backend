export IMAGE=money-keeper-backend:local-$(uuidgen)

BUILD_CONTEXT=./
DOCKERFILE=.infra/Dockerfile_dev

echo "Building image: $IMAGE"

docker build -f $DOCKERFILE $BUILD_CONTEXT -t $IMAGE
docker-compose -f .infra/development/docker-compose.yml up -d

echo "Waiting 3 secs before creating the database"
sleep 3
docker exec -it money_keeper_database mysql -u root -proot -e 'create database if not exists money_keeper';
echo "Database created (or already exists)"
