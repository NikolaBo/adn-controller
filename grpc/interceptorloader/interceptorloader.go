package interceptorloader

import (
	"plugin"

	"google.golang.org/grpc"
)

type InterceptInit interface {
	ClientInterceptors() []grpc.UnaryClientInterceptor
	ServerInterceptors() []grpc.UnaryServerInterceptor
}

func LoadInterceptors(interceptorPluginPath string) InterceptInit {
	interceptorPlugin, err := plugin.Open(interceptorPluginPath)
	if err != nil {
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
