package mux

import (
	"fmt"
	"github.com/newm4n/mihp/central/server/utils"
	"github.com/newm4n/mihp/pkg/helper"
	"net/http"
	"sort"
)

func NewMyMux() *MyMux {
	return &MyMux{
		endPoints:   make([]*Endpoint, 0),
		middlewares: make([]func(next http.Handler) http.Handler, 0),
	}
}

type MiddlewareChain struct {
	handler http.Handler
	next    http.Handler
	toCall  func(next http.Handler) http.Handler
}

type MyMux struct {
	endPoints   []*Endpoint
	middlewares []func(next http.Handler) http.Handler
}

func (m *MyMux) UseMiddleware(mw func(next http.Handler) http.Handler) {
	m.middlewares = append(m.middlewares, mw)
}

func (m *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hf, pattern := m.HandleFuncForRequest(r)
		if hf == nil {
			utils.ErrorResponse(w, "not found", http.StatusNotFound)
		} else {
			mp, _ := helper.ParsePathParams(pattern, r.URL.Path)
			if mp != nil {
				for k, v := range mp {
					// todo find an alternative way to do this. using Query() did not work
					r.URL.Query().Add(k, v)
				}
			}
			hf.ServeHTTP(w, r)
		}
	})
	if m.middlewares != nil && len(m.middlewares) > 0 {
		for i := len(m.middlewares) - 1; i >= 0; i-- {
			mw := m.middlewares[i]
			h = mw(h)
		}
	}

	h.ServeHTTP(w, r)
}

func (m *MyMux) AddRoute(pattern, method string, hFunc http.HandlerFunc) {
	ep := &Endpoint{
		PathMethod: &PathMethod{
			PathPattern: pattern,
			Method:      method,
		},
		HandleFunc: hFunc,
	}
	if len(m.endPoints) > 1 {
		sort.Slice(m.endPoints, func(i, j int) bool {
			return len(m.endPoints[j].PathMethod.String()) > len(m.endPoints[i].PathMethod.String())
		})
	}
	m.endPoints = append(m.endPoints, ep)
}

func (m *MyMux) HandleFuncForRequest(r *http.Request) (http.HandlerFunc, string) {
	for _, ep := range m.endPoints {
		if good, pattern := ep.PathMethod.MatchRequest(r); good {
			return ep.HandleFunc, pattern
		}
	}
	return nil, ""
}

type Endpoint struct {
	PathMethod *PathMethod
	HandleFunc http.HandlerFunc
}

type PathMethod struct {
	PathPattern string
	Method      string
}

func (pm *PathMethod) MatchRequest(r *http.Request) (bool, string) {
	good := helper.IsTemplateCompatible(pm.PathPattern, r.URL.Path)
	if good && r.Method == pm.Method {
		return true, pm.PathPattern
	}
	return false, ""
}

func (pm *PathMethod) String() string {
	return fmt.Sprintf("[%s]%s", pm.Method, pm.PathPattern)
}
