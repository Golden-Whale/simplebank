postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_PASSWORD=123456 -e POSTGRES_USER=root -e POSTGRES_DB=simple_bank -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover -count=1 ./...

mock:
	mockgen  -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store

proto:
	rm -f pb/*.go
	protoc --proto_path=proto \
		--go_out=pb \
		--go-grpc_out=pb \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative proto/*.proto

evans:
	evans --host localhost --port 8090 -r repl

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test mock proto evans