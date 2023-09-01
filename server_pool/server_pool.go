package serverpool

import (
	"github.com/Junkes887/load-balancer/backend"
)

type ServerPool struct {
	Backends          []*backend.Backend
	BackendsAlive     []*backend.Backend
	BackendsDontAlive []*backend.Backend
	Current           int
}

func (s *ServerPool) AddBackend(backend *backend.Backend) {
	s.Backends = append(s.Backends, backend)
	s.BackendsAlive = append(s.BackendsAlive, backend)
}

func (s *ServerPool) AddBackendAlive(backend *backend.Backend) {
	s.BackendsAlive = append(s.BackendsAlive, backend)
	index := filterBack(s.BackendsDontAlive, backend)
	s.BackendsDontAlive = append(s.BackendsDontAlive[:index], s.BackendsDontAlive[index+1:]...)

}

func (s *ServerPool) AddBackendDontAlive(backend *backend.Backend) {
	s.BackendsDontAlive = append(s.BackendsDontAlive, backend)
	index := filterBack(s.BackendsAlive, backend)
	s.BackendsAlive = append(s.BackendsAlive[:index], s.BackendsAlive[index+1:]...)
}

func filterBack(list []*backend.Backend, backend *backend.Backend) (index int) {
	for i, back := range list {
		if back.URL.String() == backend.URL.String() {
			index = i
			return
		}
	}
	index = 0
	return
}

func (s *ServerPool) GetNextPeer() *backend.Backend {
	next := s.NextIndex()
	l := len(s.BackendsAlive) + next
	for i := next; i < l; i++ {
		idx := i % len(s.BackendsAlive)
		if i == next {
			s.Current = idx
		}
		return s.BackendsAlive[idx]
	}
	return nil
}

func (s *ServerPool) NextIndex() int {
	return (s.Current + 1) % len(s.BackendsAlive)
}
