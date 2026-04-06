package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/antonioshadji/http-example/content"
	"github.com/antonioshadji/http-example/server"
)

func main() {
	handler := server.NewHandler(content.IndexHTML)

	const basePort = 8080
	const maxAttempts = 10

	for i := range maxAttempts {
		port := basePort + i
		addr := fmt.Sprintf(":%d", port)

		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Printf("port %d unavailable: %v", port, err)
			continue
		}

		log.Printf("listening on http://localhost:%d", port)
		if err := http.Serve(listener, handler); err != nil {
			log.Fatalf("server error: %v", err)
		}
		return
	}

	log.Fatalf("failed to bind to any port in range %d-%d", basePort, basePort+maxAttempts-1)
}
