package bootstrap

import (
	"fmt"

	v1 "github.com/seth16888/wxtoken/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitTokenCLient(serverAddr string) (v1.TokenClient, error) {
  grpcConn, err := grpc.NewClient(
    serverAddr,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
  )
  if err != nil {
    fmt.Println("failed to connect to server:", err)
    return nil, err
  }
  return v1.NewTokenClient(grpcConn), nil
}
