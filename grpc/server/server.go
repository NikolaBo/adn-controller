package main

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"

	"github.com/UWNetworksLab/adn-controller/grpc/interceptorloader"
	echo "github.com/UWNetworksLab/adn-controller/grpc/pb"
	"google.golang.org/grpc"
)

type server struct {
	echo.UnimplementedEchoServiceServer
}

type InterceptInit interface {
	ServerInterceptors() []grpc.UnaryServerInterceptor
}

func (s *server) Echo(ctx context.Context, x *echo.Msg) (*echo.Msg, error) {
	log.Printf("got: [%s]", x.GetBody())

	appendedBody := x.GetBody() + "1"
	msg := &echo.Msg{
		Body: appendedBody,
	}
	return msg, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interceptInit := interceptorloader.LoadInterceptors("/echoserver/severalinterceptorsplugin/severalinterceptorsplugin.so")

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptInit.ServerInterceptors()...,
		),
	)
	fmt.Printf("Starting server pod at port 9000\n")

	echo.RegisterEchoServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
