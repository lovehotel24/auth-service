package configs

import (
	"github.com/lovehotel24/booking-service/pkg/grpc/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcUserService(host string) (userpb.UserServiceClient, error) {
	gConn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return userpb.NewUserServiceClient(gConn), nil
}
