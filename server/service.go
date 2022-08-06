package server

import (
	"os"
	"path"
	"sync"
)

// ServiceEnv ...
type ServiceEnv []string

// Service ...
type Service struct {
	name string

	once sync.Once
}

// Loopkup ...
func (s *Service) Lookup(env ServiceEnv) string {
	s.once.Do(func() {
		for _, name := range env {
			v, ok := os.LookupEnv(name)
			if ok {
				s.name = v
				break
			}
		}

		if s.name == "" {
			s.name = path.Base(os.Args[0])
		}
	})

	return s.name
}
