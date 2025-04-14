package bootstrap

import (
	"fmt"

	au "github.com/seth16888/coauth/api/v1"
	px "github.com/seth16888/wxproxy/api/v1"
	ak "github.com/seth16888/wxtoken/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitTokenCLient(serverAddr string) (ak.TokenClient, error) {
  grpcConn, err := grpc.NewClient(
    serverAddr,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
  )
  if err != nil {
    fmt.Println("failed to connect to server:", err)
    return nil, err
  }
  return ak.NewTokenClient(grpcConn), nil
}

// InitAPIProxyClient initializes the gRPC client for API proxy
func InitAPIProxyClient(serverAddr string) (px.MpproxyClient, error) {
  if len(serverAddr) == 0 {
    return nil, fmt.Errorf("server address is empty")
  }
  // Create a gRPC connection to the server
  // Use insecure credentials for local development
  grpcConn, err := grpc.NewClient(
    serverAddr,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
  )
  if err != nil {
    fmt.Println("failed to connect to server:", err)
    return nil, err
  }
  return px.NewMpproxyClient(grpcConn), nil
}

// InitAuthClient initializes the gRPC client for authentication
func InitAuthClient(serverAddr string) (au.CoauthClient, error) {
  if len(serverAddr) == 0 {
    return nil, fmt.Errorf("server address is empty")
  }
  // Create a gRPC connection to the server
  // Use insecure credentials for local development
  grpcConn, err := grpc.NewClient(
    serverAddr,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
  )
  if err != nil {
    fmt.Println("failed to connect to server:", err)
    return nil, err
  }
  return au.NewCoauthClient(grpcConn), nil
}
