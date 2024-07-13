postgres:
	docker run --name postgres12 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -p 5432:5432 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	 migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	 migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
server:
	go run main.go
test:
	go test -v -cover ./...

mock:
	 mockgen -destination db/mock/store.go -package mockdb  github.com/alikhanMuslim/simpleBank/db/sqlc Store

.PHONY: createdb, dropdb, postgres, migrateup, migratedown, sqlc, test, server, mock