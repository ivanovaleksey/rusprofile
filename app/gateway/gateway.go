package gateway

import (
	"github.com/go-chi/chi"
	"net/http"
)

type Gateway struct {
	router  chi.Router
	grpcMux http.Handler
}

func New(grpcMux http.Handler) *Gateway {
	g := &Gateway{
		grpcMux: grpcMux,
	}
	g.initRouter()
	return g
}

func (g *Gateway) initRouter() {
	r := chi.NewRouter()

	r.Mount("/", g.grpcMux)
	r.Route("/doc", func(r chi.Router) {
		r.Handle("/*", http.StripPrefix("/doc", http.FileServer(http.Dir("doc"))))
	})

	g.router = r
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.router.ServeHTTP(w, r)
}
