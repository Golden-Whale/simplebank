package gapi

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/utils"
)

func (server *Server) CreateUser(c context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}
	user, err := server.store.CreateUser(c, arg)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "users_pkey", "users_email_key":
				return nil, status.Errorf(codes.AlreadyExists, err.Error())
			}
		}
		return nil, status.Errorf(codes.Internal, "cannot create user: %s", err)

	}
	rsp := &pb.CreateUserResponse{User: convertUser(user)}

	return rsp, nil
}
