package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("port", 3000, "Port number for the proxy server")
	origin := flag.String("origin", "http://localhost:3001", "Origin server URL")
	clearCache := flag.Bool("clear-cache", false, "Clear the cache before starting the server")
	flag.Parse()

	proxyServerDomain := fmt.Sprintf("localhost:%d", *port)
	proxyServer := NewProxyServer(proxyServerDomain, *origin)

	if *clearCache {
		proxyServer.ClearCache()
		fmt.Println("Cache cleaned")
		return
	}

	http.HandleFunc("/", proxyServer.HandleFunc)
	log.Printf("Starting proxy server on %s, forwarding to %s\n", proxyServerDomain, *origin)

	if err := http.ListenAndServe(proxyServerDomain, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
