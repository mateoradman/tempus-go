DB_URL=postgresql://root:secret@localhost:5432/tempus?sslmode=disable

postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root tempus

migrateup:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down

dropdb:
	docker exec -it postgres14 dropdb tempus

sqlc:
	sqlc generate

server:
	go run cmd/tempus/tempus.go

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination internal/db/mock/store.go github.com/mateoradman/tempus/internal/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

.PHONY: postgres createdb dropdb sqlc server test migrateup migratedown mock db_docs db_schema
