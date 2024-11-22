package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"ticgo/auth-service/routes"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("auth_microservices_is_up_and_running")
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go httpServer(ctx, wg)

	termSignalCh := make(chan os.Signal, 1)
	signal.Notify(termSignalCh, syscall.SIGINT, syscall.SIGTERM)
	// block untill we get notified by any interrupt or termination signal from our k8s pod
	<-termSignalCh

	// once we reached here we know that we received a signal interrupt or termination, so we need to gracefully shutdown the server, so we cancel the ctx to receive a Done signal from ctx.Done() inside the httpServer go routine
	cancel()

	// then wait untill the httpServer finish all its work and call wg.Done to gracfully close the process
	wg.Wait()

	log.Println("shutdown_gracefully_completed")
}

func httpServer(ctx context.Context, wg *sync.WaitGroup) {
	// we define the http router
	engine := gin.Default()

	// we install our routes
	healthRouter := &routes.HealthRouter{}
	healthRoutesGroup := engine.Group("/auth")
	healthRouter.InstallRouteHandlers(healthRoutesGroup)

	// define http server instance
	port := os.Getenv("PORT")
	srv := &http.Server{Addr: "0.0.0.0:" + port, Handler: engine}

	// run the server into go routine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("error_starting_http_auth_Server : ", err)
		}
	}()

	// wait for any interrupt signal from the main routine
	<-ctx.Done()

	// once we here, we need to shutdown the server, and close all open connections
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// shutdoown after 5 seconds
	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Println("error trying to shutdown the server : ", err)
	}
	wg.Done()

	log.Println("server_stopped")
}
