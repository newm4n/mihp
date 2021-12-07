package dummy

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type DummyServer struct {
	Port      int
	Alive     bool
	RandomKey string
	Srv       *http.Server
}

func init() {
	rand.Seed(time.Now().UnixMilli())
}

func (ds *DummyServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/login":
		if req.Method == http.MethodGet {
			resp.Header().Add("Content-Type", "text/plain")
			resp.Header().Add("TestToken", ds.RandomKey)
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte("OK"))
		} else {
			resp.Header().Add("Content-Type", "text/plain")
			resp.WriteHeader(http.StatusMethodNotAllowed)
			resp.Write([]byte("MethodExpr Not Allowed"))
		}
	case "/dashboard":
		for k, v := range req.Header {
			fmt.Printf("/dashboard [%s] = %s\n", k, strings.Join(v, ","))
		}
		if req.Header.Get("Authorization") == "" {
			resp.Header().Add("Content-Type", "text/plain")
			resp.WriteHeader(http.StatusUnauthorized)
			resp.Write([]byte("Unauthorized"))
		} else if req.Header.Get("Authorization") != ds.RandomKey {
			resp.Header().Add("Content-Type", "text/plain")
			resp.WriteHeader(http.StatusForbidden)
			resp.Write([]byte("Forbidden"))
		} else {
			resp.Header().Add("Content-Type", "text/plain")
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte("OK"))
		}
	default:
		resp.Header().Add("Content-Type", "text/plain")
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte("Not Found"))
	}
}

func (ds *DummyServer) Stop() {
	if ds.Alive {
		_ = ds.Srv.Shutdown(context.Background())
		ds.Alive = false
	}
}

func (ds *DummyServer) Start() {
	ds.Port = rand.Intn(50000) + 10000

	if !ds.Alive {
		ds.Srv = &http.Server{
			Addr:              fmt.Sprintf("0.0.0.0:%d", ds.Port),
			Handler:           ds,
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       10 * time.Second,
		}
		fmt.Printf("Dummy server alive at %d\n", ds.Port)
		byteArr := make([]byte, 20)
		for i := 0; i < 20; i++ {
			byteArr[i] = byte(rand.Intn(25) + 65)
		}
		ds.RandomKey = string(byteArr)
		go func() {
			_ = ds.Srv.ListenAndServe()
		}()
		ds.Alive = true
	}
}
