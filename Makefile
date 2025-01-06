dev:
	docker-compose up go-server-dev vue-client-dev --remove-orphans
build:
	docker-compose build
go-test:
	docker-compose run --remove-orphans go-server-dev scripts/test.sh
