# gRPC echo server

This is a simple Echo server built using Go and gRPC, with a set of example interceptors that are loaded via [Go plugin](https://pkg.go.dev/plugin).

## Run as docker container
To run the server as a Docker container, follow these steps:
- `docker build --tag echo-frontend -f Dockerfile-frontend .`
- `docker build --tag echo-server -f Dockerfile-server  .`
- `docker network create test`
- `docker run --rm -d --net test -p 9000:9000 --name server echo-server`
- `docker run --rm -d --net test -p 8080:8080 --name frontend echo-frontend`
- `curl http://localhost:8080/echo`

To run the server with hotswap frontend, follow these steps:
- `docker build --tag echo-hotswapfrontend -f Dockerfile-hotswapfrontend .`
- `docker build --tag echo-server -f Dockerfile-server  .`
- `docker network create test`
- `docker run --rm -d --net test -p 9000:9000 --name server echo-server`
- `docker run --rm -d --net test -p 8080:8080 --name frontend echo-hotswapfrontend`
- `curl http://localhost:8080/echo`

You can then update the interceptors without application restart:
- `go build  -C ./singleinterceptorplugin/ -buildmode=plugin  .`
- `docker cp ./singleinterceptorplugin/singleinterceptorplugin.so frontend:/interceptors.so1`
- `curl http://localhost:8080/echo`
- `go build  -C ./severalinterceptorsplugin/ -buildmode=plugin  .`
- `docker cp ./singleinterceptorplugin/severalinterceptorsplugin.so frontend:/interceptors.so2`
- `curl http://localhost:8080/echo`

## Note
Work in progress. Potential issues:
- Application and plugins must be built with the same version of Go and any libraries shared by plugin and application must be the exact same source.
- May not be able to remove plugins once they are loaded.
- Go tries to prevent loading the same plugin twice (https://github.com/golang/go/issues/47298) so when loading multiple plugins they must be built as separate packages.