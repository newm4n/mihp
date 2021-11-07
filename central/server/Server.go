package server

import (
	"context"
	"fmt"
	mux "github.com/hyperjumptech/hyper-mux"
	"github.com/newm4n/mihp/central/server/handlers"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer struct {
	Host      string
	Port      int64
	TheServer *http.Server
	Mux       *mux.HyperMux
}

func (s *HttpServer) Start() {
	var wait time.Duration

	// StartUpTime records first ime up
	startUpTime := time.Now()

	// Initializing muxes and routing
	s.Mux = mux.NewHyperMux()
	s.Mux.UseMiddleware(mux.NewCORSMiddleware(mux.DefaultCORSOption))
	s.Mux.UseMiddleware(mux.ContextSetterMiddleware)
	s.Mux.UseMiddleware(mux.GZIPCompressMiddleware)
	s.Mux.UseMiddleware(handlers.BearerCheckMiddleware)
	handlers.Routing(s.Mux)

	// Initializing servers etc

	defer s.Shutdown()

	if s.TheServer == nil {
		s.TheServer = &http.Server{
			Addr:              fmt.Sprintf("%s:%d", s.Host, s.Port),
			Handler:           s.Mux,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       30 * time.Second,
		}
	}
	go func() {
		err := s.TheServer.ListenAndServe()
		if err != nil {
			logrus.Error(err.Error())
		}
	}()

	gracefulStop := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(gracefulStop, os.Interrupt)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	// Block until we receive our signal.
	<-gracefulStop

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	s.TheServer.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logrus.Info("shutting down........ bye")

	t := time.Now()
	upTime := t.Sub(startUpTime)
	fmt.Println("server was up for : ", upTime.String(), " *******")
	os.Exit(0)
}

func (s *HttpServer) Shutdown() {
	// Do all clean-up here

}
