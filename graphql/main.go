//go:generate go run github.com/99designs/gqlgen
package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	router.Handle("/graphql", handler.GraphQL(s.ToExecutableSchema()))
	router.Handle("/playground", handler.Playground("Spidey", "/graphql"))
	log.Fatal(http.ListenAndServe(":8000", router))
}
