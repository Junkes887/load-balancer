package healthcheck

import (
	"log"
	"net"
	"net/url"
	"time"

	serverpool "github.com/Junkes887/load-balancer/server_pool"
)

func LauchHealthCheck(serverPool *serverpool.ServerPool) {
	t := time.NewTicker(time.Minute * 5)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			healthCheck(serverPool)
			log.Println("Health check completed")
		}
	}
}

func healthCheck(serverPool *serverpool.ServerPool) {
	for _, back := range serverPool.BackendsAlive {
		status := "up"
		alive := isBackendAlive(back.URL)
		if !alive {
			status = "down"
			serverPool.AddBackendDontAlive(back)
		}
		log.Printf("%s [%s]\n", back.URL, status)
	}
	for _, back := range serverPool.BackendsDontAlive {
		status := "down"
		alive := isBackendAlive(back.URL)
		if alive {
			status = "up"
			serverPool.AddBackendAlive(back)
		}
		log.Printf("%s [%s]\n", back.URL, status)
	}
}

func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	defer conn.Close()
	return true
}
