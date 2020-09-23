//go:generate go run github.com/99designs/gqlgen
package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/cors"
)

type appConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
	RedisURL   string `envconfig:"REDIS_URL"`
}

func main() {
	var cfg appConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL, cfg.RedisURL)
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	router.Handle("/graphql", handler.GraphQL(s.ToExecutableSchema(), handler.WebsocketUpgrader(websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})))
	router.Handle("/playground", handler.Playground("Spidey", "/graphql"))
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "DELETE", "POST", "PUT"},
	}).Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}
