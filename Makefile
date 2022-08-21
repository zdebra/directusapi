test: start_server unit stop_server

start_server:
	docker-compose pull
	docker-compose up -d
	docker-compose run --rm directus install --email email@example.com --password d1r3ctu5

stop_server:
	docker-compose down

unit:
	go test -count=1 ./...


