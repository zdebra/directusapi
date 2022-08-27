DOCKER_COMPOSE_DIRECTUS_V8 := docker-compose-directus-8.yml
DIRECTUS_ADMIN_EMAIL := email@example.com
DIRECTUS_ADMIN_PASSWORD := d1r3ctu5

e2e_test: start_v8_server go_test stop_v8_server

start_v8_server:
	docker-compose --file $(DOCKER_COMPOSE_DIRECTUS_V8) pull
	docker-compose --file $(DOCKER_COMPOSE_DIRECTUS_V8) up -d
	docker-compose --file $(DOCKER_COMPOSE_DIRECTUS_V8) run --rm directus install --email $(DIRECTUS_ADMIN_EMAIL) --password $(DIRECTUS_ADMIN_PASSWORD)

stop_v8_server:
	docker-compose --file $(DOCKER_COMPOSE_DIRECTUS_V8) down

go_test:
	go test -count=1 ./...


