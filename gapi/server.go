package gapi

import (
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/token"
	"simplebank/utils"
)

// Server serves GRPC requests for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new GRPC server and setup routing
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenmaker, err := token.NewPasetoMaker(config.TokenSysmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenmaker, config: config,
	}

	return server, nil
}
