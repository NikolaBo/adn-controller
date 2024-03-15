package main

import (
	"fmt"
	"log"
	"net/http"
	"plugin"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	// "github.com/UWNetworksLab/adn-controller/grpc/interceptors/null"

	echo "github.com/UWNetworksLab/adn-controller/grpc/pb"
)

type InterceptInit interface {
	CacheClient() grpc.UnaryClientInterceptor
}

func initHandler() (func(writer http.ResponseWriter, request *http.Request), func()) {
	// nullOpts := []null.CallOption{null.WithMessage("Null"),}
	// aclOpts := []acl.CallOption{acl.WithContent("client")}

	interceptorPlugin, err := plugin.Open("/echoserver/interceptors/cache/cache.so")
	if err != nil {
		panic("error loading interceptor plugin so")
	}

	symInterceptInit, err := interceptorPlugin.Lookup("InterceptInit")
	if err != nil {
		panic("error locating interceptor in plugin so")
	}

	var interceptInit InterceptInit
	interceptInit, ok := symInterceptInit.(InterceptInit)
	if !ok {
		panic("error casting interceptInit")
	}

	conn, err := grpc.Dial(
		"server:9000",
		grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(
			// null.NullClient(nullOpts...),
			// rate.RateClient(),
			// acl.ACLClient(aclOpts...),
			// mutate.MutateClient(),
			// cache.CacheClient(),
			interceptInit.CacheClient(),
		),
	)
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
			request_body := strings.Replace(request.URL.String(), "/", "", -1)
			fmt.Printf("Got request with body: %s\n", request_body)

			c := echo.NewEchoServiceClient(conn)

			message := echo.Msg{
				Body: request_body,
			}

			response, err := c.Echo(context.Background(), &message)
			if err != nil {
				fmt.Fprintf(writer, "Echo server returns an error.\n")
				log.Printf("Error when calling echo: %s", err)
			} else {
				fmt.Fprintf(writer, "Echo request finished! Length of the request is %d\n", len(response.Body))
				log.Printf("Response from server: %s", response.Body)
			}
		}, func() {
			conn.Close()
		}
}

func main() {
	handler, cleanUp := initHandler()
	defer cleanUp()

	http.HandleFunc("/", handler)

	fmt.Printf("Starting frontend pod at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
