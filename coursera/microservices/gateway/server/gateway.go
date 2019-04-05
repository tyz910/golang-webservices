package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	// "../session"
	"coursera/microservices/gateway/session"
)

func main() {
	proxyAddr := ":8080"
	serviceAddr := "127.0.0.1:8081"

	go gRPCService(serviceAddr)
	HTTPProxy(proxyAddr, serviceAddr)
}

/*
curl -X POST -k http://localhost:8080/v1/session/create -H "Content-Type: text/plain" -d '{"login":"rvasily", "useragent": "chrome"}'
curl http://localhost:8080/v1/session/check/XVlBzgbaiC
curl -X POST -k http://localhost:8080/v1/session/delete -H "Content-Type: text/plain" -d '{"ID":"XVlBzgbaiC"}'
*/

func gRPCService(serviceAddr string) {
	lis, err := net.Listen("tcp", serviceAddr)
	if err != nil {
		log.Fatalln("failed to listen TCP port", err)
	}

	server := grpc.NewServer()

	session.RegisterAuthCheckerServer(server, NewSessionManager())

	fmt.Println("starting gRPC server at " + serviceAddr)
	server.Serve(lis)
}

func HTTPProxy(proxyAddr, serviceAddr string) {
	grcpConn, err := grpc.Dial(
		serviceAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("failed to connect to grpc", err)
	}
	defer grcpConn.Close()

	grpcGWMux := runtime.NewServeMux()

	err = session.RegisterAuthCheckerHandler(
		context.Background(),
		grpcGWMux,
		grcpConn,
	)
	if err != nil {
		log.Fatalln("failed to start HTTP server", err)
	}

	mux := http.NewServeMux()
	// отправляем в прокси только то что нужно
	mux.Handle("/v1/session/", grpcGWMux)

	mux.HandleFunc("/", helloWorld)

	fmt.Println("starting HTTP server at " + proxyAddr)
	log.Fatal(http.ListenAndServe(proxyAddr, mux))

}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "URL:", r.URL.String())
}
