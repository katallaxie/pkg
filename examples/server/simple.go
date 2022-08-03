package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/katallaxie/pkg/debug"
	"github.com/katallaxie/pkg/server"
)

type srv struct {
}

func (s *srv) Start(ctx context.Context, ready server.ReadyFunc) func() error {
	return func() error {
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome"))
		})

		ready()

		s := http.Server{Addr: ":3000", Handler: chi.ServerBaseContext(ctx, r)}

		if err := s.ListenAndServe(); err != nil {
			return err
		}

		return nil
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, _ := server.WithContext(ctx)

	d := debug.New(
		debug.WithPprof(),
		debug.WithStatusAddr(":8443"),
	)
	s.Listen(d, false)
	s.Listen(&srv{}, true)

	err := s.Wait()

	var e *server.Error
	if errors.As(err, &e) {
		fmt.Println(e.Err)
		os.Exit(1)
	}
}
