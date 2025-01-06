dev:
	docker-compose up go-server-dev vue-client-dev --remove-orphans
build:
	docker-compose build
go-test:
	docker-compose run --rm go-server-dev scripts/test.sh
