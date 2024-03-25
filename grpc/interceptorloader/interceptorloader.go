package interceptorloader

import (
	"fmt"
	"plugin"

	"google.golang.org/grpc"
)

type InterceptInit interface {
	ClientInterceptors() []grpc.UnaryClientInterceptor
	ServerInterceptors() []grpc.UnaryServerInterceptor
}

func LoadInterceptors(interceptorPluginPath string) InterceptInit {
	// TODO: return err instead of panicking
	interceptorPlugin, err := plugin.Open(interceptorPluginPath)
	if err != nil {
		fmt.Printf("loading error: %v\n", err)
		panic("error loading interceptor plugin so")
	}

	symInterceptInit, err := interceptorPlugin.Lookup("InterceptInit")
	if err != nil {
		panic("error locating interceptor in plugin so")
	}

	interceptInit, ok := symInterceptInit.(InterceptInit)
	if !ok {
		panic("error casting interceptInit")
	}

	return interceptInit
}
