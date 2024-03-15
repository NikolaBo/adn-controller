package rate

import (
	"log"
	"sync"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RateClient() grpc.UnaryClientInterceptor {
	limit := uint(20)
	perSec := 0.1

	lastTs := time.Now()
	var mu sync.Mutex
	token := uint(0)

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Println("Running RateUnaryClientInterceptor")

		now := time.Now()
		mu.Lock()

		tokenToAdd := uint((float64((now.Second())) - float64(lastTs.Second())) * perSec)
		log.Printf("tokenToAdd is: %v", tokenToAdd)
		log.Printf("tokenToAdd + currentToken is: %v", tokenToAdd+token)
		token = min(tokenToAdd+token, limit)
		log.Printf("currentToken is: %v", token)
		lastTs = now

		if token < 1 {
			mu.Unlock()
			return status.Error(codes.Aborted, "request blocked by Rate Limiter.")
		}
		token--
		mu.Unlock()

		err := invoker(ctx, method, req, reply, cc, opts...)

		return err
	}
}
