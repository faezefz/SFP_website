postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

migrateup:
	migrate path db/migration -database "postgres://root:secret@postgres12:5432/sfp_db?ssl_mode=disable" -verbose up

migratedown:
	migrate path db/migration -database "postgres://root:secret@localhost:5432/sfp_db?ssl_mode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root sfp_db

dropdb:
	docker exec -it postgres12 dropdb sfp_db

.PHONY: postgres migrateup migratedown createdb dropdb sqlc test