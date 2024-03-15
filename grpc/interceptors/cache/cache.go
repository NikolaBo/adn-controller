package cache

import (
	"fmt"
	"log"
	"sync"

	echo "github.com/UWNetworksLab/adn-controller/grpc/pb"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

func CacheClient() grpc.UnaryClientInterceptor {
	var request_bodies sync.Map

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Println("Running CacheUnaryClientInterceptor")

		m, ok := req.(*echo.Msg)
		if !ok {
			panic("Incorrect message type.")
		}

		cachedResult, cacheHit := request_bodies.Load(m.GetBody())
		if cacheHit {
			fmt.Printf("Cache hit!!! body: %s\n", m.GetBody())
			*reply.(*echo.Msg) = echo.Msg{
				Body: cachedResult.(string),
			}
			return nil
		} else {
			fmt.Printf("Cache miss!!! body: %s\n", m.GetBody())
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		responseBody := reply.(*echo.Msg).GetBody()
		fmt.Printf("Inserting request to cache. Body : %s\n", responseBody)
		request_bodies.Store(m.GetBody(), responseBody)

		return err
	}
}
