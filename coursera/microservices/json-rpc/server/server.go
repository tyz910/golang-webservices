package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

/*
{
   "jsonrpc":"2.0",
   "id":1,
   "method":"SessionManager.Create",
   "params":[
      {
         "login":"rvasily",
         "useragent":"chrome"
      }
   ]
}
*/

/*

curl -v -X POST -H "Content-Type: application/json" -H "X-Auth: 123" -d '{"jsonrpc":"2.0", "id": 1, "method": "SessionManager.Create", "params": [{"login":"rvasily", "useragent": "chrome"}]}' http://localhost:8081/rpc

curl -v -X POST -H "Content-Type: application/json" -H "X-Auth: 123" -d '{"jsonrpc":"2.0", "id": 2, "method": "SessionManager.Check", "params": [{"id":"XVlBzgbaiC"}]}' http://localhost:8081/rpc

*/

type Handler struct {
	rpcServer *rpc.Server
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("rpc auth: ", r.Header.Get("X-Auth"))

	serverCodec := jsonrpc.NewServerCodec(&HttpConn{
		in:  r.Body,
		out: w,
	})
	w.Header().Set("Content-type", "application/json")
	err := h.rpcServer.ServeRequest(serverCodec)
	if err != nil {
		log.Printf("Error while serving JSON request: %v", err)
		http.Error(w, `{"error":"cant serve request"}`, 500)
	} else {
		w.WriteHeader(200)
	}
}

func main() {
	sessManager := NewSessManager()

	server := rpc.NewServer()
	server.Register(sessManager)

	sessionHandler := &Handler{
		rpcServer: server,
	}
	http.Handle("/rpc", sessionHandler)

	fmt.Println("starting server at :8081")
	http.ListenAndServe(":8081", nil)

}
