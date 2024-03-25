package main

import (
	"github.com/UWNetworksLab/adn-controller/grpc/interceptors/fault"
	"google.golang.org/grpc"
)

type interceptInit struct{}

func (interceptInit) ClientInterceptors() []grpc.UnaryClientInterceptor {
	return []grpc.UnaryClientInterceptor{fault.FaultClient()}
}

func (interceptInit) ServerInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{fault.FaultServer()}
}

var InterceptInit interceptInit
