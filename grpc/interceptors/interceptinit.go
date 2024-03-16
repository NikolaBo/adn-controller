package main

import (
	"github.com/UWNetworksLab/adn-controller/grpc/interceptors/acl"
	"github.com/UWNetworksLab/adn-controller/grpc/interceptors/cache"
	"github.com/UWNetworksLab/adn-controller/grpc/interceptors/fault"
	"github.com/UWNetworksLab/adn-controller/grpc/interceptors/mutate"
	"google.golang.org/grpc"
)

type interceptInit struct{}

func (interceptInit) ClientInterceptors() []grpc.UnaryClientInterceptor {
	return []grpc.UnaryClientInterceptor{cache.CacheClient(), mutate.MutateClient(), fault.FaultClient()}
}

func (interceptInit) ServerInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{acl.ACLServer(), mutate.MutateServer(), fault.FaultServer()}
}

var InterceptInit interceptInit
