package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/elisalimli/serverless-url-alias/domain"
	"github.com/elisalimli/serverless-url-alias/internal/api"
	"github.com/elisalimli/serverless-url-alias/internal/initializers"
	"github.com/elisalimli/serverless-url-alias/internal/sheets"
)

func init() {

	initializers.Initialize()
}

func main() {
	log.Println("-------------------------")
	ctx := context.Background()

	client, err := sheets.NewClient(ctx)
	if err != nil {
		log.Fatalf("error occured while creating new service : %v", err)
	}
	d := domain.NewDomain(*client)
	handlers := &api.Handlers{Domain: d}

	http.HandleFunc("/", handlers.Redirect)
	listenAddr := net.JoinHostPort(initializers.Host, initializers.Port)

	log.Printf("starting to listen at %v", listenAddr)
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
