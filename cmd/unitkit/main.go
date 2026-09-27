package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/relentlessworks/unitkit/internal/api"
	"github.com/relentlessworks/unitkit/internal/config"
)

func main() {
	cfg := config.Load()

	handler := api.New()

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: handler,
	}

	fmt.Fprintf(os.Stderr, "unitkit listening on %s\n", cfg.Addr)
	log.Fatal(server.ListenAndServe())
}
