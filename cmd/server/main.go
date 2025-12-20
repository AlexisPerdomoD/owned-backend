package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

// start point baby
func main() {
	l := slog.New(slog.NewTextHandler(os.Stdout, nil))
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello Ownned"))
		if err != nil {
			l.Warn("some error happend", "err", err)
		}
	})

	PORT := 9090
	l.Info("server starting at:", "port", PORT)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), r); err != nil && err != http.ErrServerClosed {
		l.Error("server could not start", "err", err)
	}

}
