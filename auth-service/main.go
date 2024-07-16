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
	log.Println("auth microservices is up and running")
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go httpServer(ctx, wg)

	termSignalCh := make(chan os.Signal, 1)
	signal.Notify(termSignalCh, syscall.SIGINT, syscall.SIGTERM)
	// block untill we get notified by any interrupt or termination signal
	<-termSignalCh

	// once we reached here we know that we received a signal interrupt or termination, so we need to gracefully shutdown the server, so we cancel the ctx to receive a Done msg from ctx.Done() inside the httpServer method
	cancel()

	// then wait untill the httpServer finish all its work and call wg.Done to gracfully close the process
	wg.Wait()

	log.Println("shutdown gracefully completed")
}

func httpServer(ctx context.Context, wg *sync.WaitGroup) {
	engine := gin.Default()

	healthRouter := &routes.HealthRouter{}
	healthRoutesGroup := engine.Group("/auth")
	healthRouter.InstallRouteHandlers(healthRoutesGroup)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}

	srv := &http.Server{Addr: port, Handler: engine}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("error_starting_http_auth_Server : ", err)
		}
	}()

	select {
	case <-ctx.Done():
		defer wg.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// shutdoown after 5 seconds
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Println("error trying to shutdown the server : ", err)
		}
	}
	log.Println("server stopped")

}
