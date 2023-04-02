package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/katallaxie/pkg/debug"
	"github.com/katallaxie/pkg/server"
)

type srv struct{}

func (s *srv) Start(ctx context.Context, ready server.ReadyFunc, run server.RunFunc) func() error {
	return func() error {
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome"))
		})

		time.Sleep(3 * time.Second)

		run(func() error {
			ticker := time.NewTicker(1 * time.Second)

			for {
				select {
				case <-ticker.C:
					log.Print(time.Now())
				case <-ctx.Done():
					ticker.Stop()
					return nil
				}
			}
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
	s.SetLimit(2)

	s.Listen(&srv{}, true)
	d := debug.New(
		debug.WithPprof(),
		debug.WithStatusAddr(":8888"),
	)
	s.Listen(d, false)

	log.Printf("starting %s", server.Service.Lookup(server.ServiceEnv{}))
	if err := s.Wait(); errors.Is(err, &server.Error{}) {
		fmt.Println(err)
		os.Exit(1)
	}
}
