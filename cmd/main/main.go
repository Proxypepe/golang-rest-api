package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"restapi/internal/user"
	"restapi/pkg/logging"
	"time"
)


func main() {
	logger := logging.GetLogger()
	logger.Info("Create router")
	router := httprouter.New()
	logger.Info("Register user handler")
	handler := user.NewHandler(logger)
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
