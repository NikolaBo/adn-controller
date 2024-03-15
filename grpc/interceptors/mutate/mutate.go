package mutate

import (
	"log"
	"strings"

	echo "github.com/UWNetworksLab/adn-controller/grpc/pb"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

func MutateClient() grpc.UnaryClientInterceptor {

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Println("Running MutateUnaryClientInterceptor")

		if m, ok := req.(*echo.Msg); ok {
			m.Body = strings.Replace(m.Body, "secret", "modified", -1)
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

func MutateServer() grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Println("Running MutateUnaryServerInterceptor")

		reply, err := handler(ctx, req)
		if m, ok := reply.(*echo.Msg); ok {
			m.Body = strings.Replace(m.Body, "hidden", "modified", -1)
		}

		return reply, err
	}
}
