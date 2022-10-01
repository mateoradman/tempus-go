postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root tempus

migrateup:
	migrate -path internal/db/migration -database "postgresql://root:secret@localhost:5432/tempus?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/db/migration -database "postgresql://root:secret@localhost:5432/tempus?sslmode=disable" -verbose down

dropdb:
	docker exec -it postgres14 dropdb tempus

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb sqlc server test