DC := docker-compose -p storage-gateway -f docker/docker-compose.yml

run:
	go run cmd/storage-gateway/main.go -conf config/config.local.json

# Stops and removes all Docker containers, networks, and volumes
docker-clean:
	$(DC) down --remove-orphans --volumes

# Cleans up containers and then starts up all the services defined in docker-compose
docker-up:
	make docker-clean
	$(DC) up --force-recreate -d

# Cleans up containers and then starts up only object-storage nodes
docker-storage:
	make docker-clean
	$(DC) up --force-recreate -d amazin-object-storage-node-1 amazin-object-storage-node-2 amazin-object-storage-node-3

lint:
	staticcheck ./...