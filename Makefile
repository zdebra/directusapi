e2e_test: start_server go_test stop_server

start_server:
	docker-compose pull
	docker-compose up -d
	docker-compose run --rm directus install --email email@example.com --password d1r3ctu5

stop_server:
	docker-compose down

go_test:
	go test -count=1 ./...


