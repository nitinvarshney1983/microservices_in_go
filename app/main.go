package main

import (
	"context"
	"log"
	"microservices/app/configurations"
	"microservices/app/handlers"
	"microservices/app/logging"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	loadBanner()
	configArgs := configurations.ConfigArgs{
		Name: "appConfigs",
		Path: []string{"./configs", "${HOME}/configs"},
	}
	configurations.Setup(&configArgs)
}

func main() {

	/*
		serveMux is a httpHandler. If you don't specify an http handler with http it uses defaultServeMux
		ServeMux, basically maps the path to a particular function
	*/
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	logging.SettingUpLoggerWithconfigurations(log, configurations.Get("LOGGING"))
	//defining handlers
	//rootHandler := handlers.NewRootHandler(log)
	//goodbyeHandler := handlers.NewGoodBuy(log)
	productHandler := handlers.NewProduct(log)

	//created a new ServeMux which is a multiplexer. It will check the routs and maps them to handlers
	sm := http.NewServeMux()
	//sm.Handle("/", rootHandler)
	//sm.Handle("/goodbye", goodbyeHandler)
	//sm.Handle("/products", productHandler)
	sm.Handle("/", productHandler)

	//defining new Server
	server := &http.Server{
		Addr:         ":9090",
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      sm,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Some error occurred at server level %v", err)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	sig := <-sigChan
	log.Panicln("Received terminate, graceful shutdown", sig)

	//for graceful shutdown
	tc, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	server.Shutdown(tc)
}

func loadBanner() {
	if data, err := os.ReadFile("./configs/banner.txt"); err == nil {
		log.Println(string(data))
		log.Println("")
	}
}
