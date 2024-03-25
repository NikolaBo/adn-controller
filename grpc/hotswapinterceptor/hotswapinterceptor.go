package hotswapinterceptor

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/UWNetworksLab/adn-controller/grpc/interceptorloader"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

func ClientInterceptor(pluginPathPrefix string) grpc.UnaryClientInterceptor {
	var currentChain grpc.UnaryClientInterceptor
	var lastName string

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		fileName := chooseFile(pluginPathPrefix)

		if lastName != fileName && fileName != "" {
			lastName = fileName

			interceptInit := interceptorloader.LoadInterceptors(fileName)
			interceptors := interceptInit.ClientInterceptors()
			log.Println("loaded new client interceptors")
			log.Printf("interceptors: %v\n", interceptors)

			if len(interceptors) == 0 {
				currentChain = nil
			} else if len(interceptors) == 1 {
				currentChain = interceptors[0]
			} else {
				currentChain = func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
					return interceptors[0](ctx, method, req, reply, cc, getChainUnaryInvoker(interceptors, 0, invoker), opts...)
				}
			}
		}

		if currentChain == nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		return currentChain(ctx, method, req, reply, cc, invoker, opts...)
	}
}

func chooseFile(prefix string) string {
	var highestNumber int
	var highestFile string

	dir, prefix := filepath.Split(prefix)

	files, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) {
			suffix := strings.TrimPrefix(file.Name(), prefix)
			if num, err := strconv.Atoi(suffix); err == nil {
				if num > highestNumber {
					highestNumber = num
					highestFile = file.Name()
				}
			} else if highestFile == "" {
				highestFile = file.Name()
			}
		}
	}

	if highestFile == "" {
		return ""
	}
	return dir + highestFile
}

// chained interceptor generation from https://github.com/grpc/grpc-go/blob/55cd7a68b3c18a0f76ea9c1be37221a5b901a798/clientconn.go#L435
func getChainUnaryInvoker(interceptors []grpc.UnaryClientInterceptor, curr int, finalInvoker grpc.UnaryInvoker) grpc.UnaryInvoker {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return interceptors[curr+1](ctx, method, req, reply, cc, getChainUnaryInvoker(interceptors, curr+1, finalInvoker), opts...)
	}
}
