DB_URL=postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable

network:
	docker network create bank-network

postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_PASSWORD=123456 -e POSTGRES_USER=root -e POSTGRES_DB=simple_bank -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "${DB_URL}" -verbose up

migratedown:
	migrate -path db/migration -database "${DB_URL}"  -verbose down

migrateup1:
	migrate -path db/migration -database "${DB_URL}"  -verbose up 1

migratedown1:
	migrate -path db/migration -database "${DB_URL}"  -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover -count=1 ./...

server:
	go run main.go

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

.PHONY:network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock proto evans