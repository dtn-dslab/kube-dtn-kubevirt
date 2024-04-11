package common

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectDaemon(ctx context.Context, ip string) (*grpc.ClientConn, error) {
	daemonAddr := fmt.Sprintf("passthrough:///%s:%d", ip, DefaultPort)
	conn, err := grpc.Dial(daemonAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
