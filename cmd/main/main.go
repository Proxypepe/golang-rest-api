package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"restapi/internal/config"
	"restapi/internal/user"
	"restapi/pkg/logging"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	//cfgMongo := cfg.MongoDB

	//mongoDBClient, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username,
	//	cfgMongo.Password, cfgMongo.Database, "")
	//if err != nil {
	//	panic(err)
	//}

	//storage := db.NewStorage(mongoDBClient, cfgMongo.Collection, logger)
	//
	//user1 := user.User{
	//	ID:           "",
	//	Username:     "lala",
	//	PasswordHash: "sdsdadad",
	//	Email:        "dsdsd",
	//}

	//userId, err := storage.Create(context.Background(), user1)
	//if err != nil {
	//	panic(err)
	//}
	//logger.Info(userId)

	logger.Info("Register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)
	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("Start app")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("Create socket")
		socketPath := path.Join(appDir, "app.sock")
		logger.Debugf("socket path %s", socketPath)

		logger.Info("Create unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
	} else {
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:           router,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	log.Fatalln(server.Serve(listener))
}
