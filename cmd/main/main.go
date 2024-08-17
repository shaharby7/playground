package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/julienschmidt/httprouter"
	"github.com/shaharby7/playground/pkg/greeter"

	"github.com/shaharby7/playground/pkg/util"
	"github.com/swaggest/swgui/v5emb"

	"io"
	"net"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	router := httprouter.New()
	port := 3000
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	router.Handle(
		"POST",
		"/api/greet",
		func(writer http.ResponseWriter, request *http.Request, p httprouter.Params) {
			body, err := io.ReadAll(request.Body)
			if err != nil {
				writer.WriteHeader(500)
				io.WriteString(writer, "something went wrong: body not found")
				return
			}
			in := new(greeter.GreetInput)
			if err := json.Unmarshal([]byte(body), in); err != nil {
				writer.WriteHeader(400)
				io.WriteString(writer, "something went wrong: json could not be parsed")
				return
			}
			out, err := greeter.Greet(context.TODO(), in)
			if err != nil {
				writer.WriteHeader(400)
				io.WriteString(writer, "something went wrong: Greet can't be created")
				return
			}
			resp, err := json.Marshal(out)
			if err != nil {
				writer.WriteHeader(400)
				io.WriteString(writer, "something went wrong: Greet can't be marshaled")
				return
			}
			writer.Header().Add("Content-Type", "application/json")
			io.WriteString(writer, string(resp))
		},
	)
	router.ServeFiles("/api/*filepath", http.Dir("api/"))
	router.Handle(
		"GET",
		"/swagger/*any",
		util.WrapHandler(v5emb.NewHandler(
			"Greeter",
			"http://localhost:3000/api/openapi.json",
			"/swagger/index.html",
		)),
	)
	wg.Add(1)
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		panic(err.Error())
	}
	go func() {
		if err := server.Serve(listener); err != nil {
			panic(err)
		}
	}()
	fmt.Printf("http server is listening on port: %d", port)
	wg.Wait()
}
