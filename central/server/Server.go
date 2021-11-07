package server

import (
	"context"
	"fmt"
	mux "github.com/hyperjumptech/hyper-mux"
	"github.com/newm4n/mihp/central/server/handlers"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

type HttpServer struct {
	Host        string
	Port        int64
	alive       bool
	TheServer   *http.Server
	Middlewares []fasthttp.RequestHandler
	Mux         *mux.HyperMux
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (s *HttpServer) Start() error {
	if !s.alive {
		s.Mux = mux.NewHyperMux()
		s.Mux.UseMiddleware(mux.NewCORSMiddleware(mux.DefaultCORSOption))
		s.Mux.UseMiddleware(mux.ContextSetterMiddleware)
		s.Mux.UseMiddleware(handlers.BearerCheckMiddleware)
		handlers.Routing(s.Mux)
		if s.TheServer == nil {
			s.TheServer = &http.Server{
				Addr:              fmt.Sprintf("%s:%d", s.Host, s.Port),
				Handler:           s,
				ReadHeaderTimeout: 10 * time.Second,
				ReadTimeout:       10 * time.Second,
				WriteTimeout:      10 * time.Second,
				IdleTimeout:       30 * time.Second,
			}
		}
		go func() {
			err := s.TheServer.ListenAndServe()
			if err != nil {
				panic(err.Error())
			} else {
				s.alive = true
			}
		}()
	} else {
		return fmt.Errorf("server already alive")
	}
	return nil
}

func (s *HttpServer) Shutdown() error {
	if s.alive {
		err := s.TheServer.Shutdown(context.Background())
		if err != nil {
			return err
		}
		s.alive = false
	} else {
		return fmt.Errorf("server already down")
	}
	return nil
}

type LogrusFastHttpLogger struct {
	LogLevel logrus.Level
	Logger   *logrus.Entry
}

func (ll *LogrusFastHttpLogger) Printf(format string, args ...interface{}) {
	switch ll.LogLevel {
	case logrus.DebugLevel:
		ll.Logger.Debugf(format, args)
	}
}
