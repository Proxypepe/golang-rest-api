package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"restapi/internal/user"
	"time"
)


func main() {
	log.Println("Create router")
	router := httprouter.New()
	log.Println("Register user handler")
	handler := user.NewHandler()
	handler.Register(router)
	start(router)
}

func start(router *httprouter.Router) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler: router,
		WriteTimeout: 15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	log.Fatalln(server.Serve(listener))
}
