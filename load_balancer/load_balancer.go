package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Junkes887/load-balancer/backend"
	serverpool "github.com/Junkes887/load-balancer/server_pool"
)

type LoadBalancer struct {
	SP *serverpool.ServerPool
}

func (lb LoadBalancer) Proxy(w http.ResponseWriter, r *http.Request) {
	peer := lb.SP.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func (lb LoadBalancer) AddBackends(urls []string) {
	for _, u := range urls {
		serverUrl, err := url.Parse(u)
		if err != nil {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", serverUrl.Host, e.Error())

			lb.dontAlive(serverUrl)
		}

		lb.SP.AddBackend(&backend.Backend{
			URL:          serverUrl,
			ReverseProxy: proxy,
		})
		log.Printf("Configured server: %s\n", serverUrl)
	}
}

func (lb LoadBalancer) dontAlive(backendUrl *url.URL) {
	for _, b := range lb.SP.BackendsAlive {
		if b.URL.String() == backendUrl.String() {
			lb.SP.AddBackendDontAlive(b)
			break
		}
	}
}
