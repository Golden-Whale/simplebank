package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"simplebank/api"
	db "simplebank/db/sqlc"
)

const (
	dbSource = "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable"
	address  = "127.0.0.1:8080"
)

func main() {
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(connPool)

	server := api.NewServer(store)
	os.IsNotExist(server.Start(address))
}
