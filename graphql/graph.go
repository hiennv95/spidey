package main

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/go-redis/redis"
	"github.com/tinrab/spidey/account"
	"github.com/tinrab/spidey/catalog"
	"github.com/tinrab/spidey/order"
)

// Server is server
type Server struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
	redisClient   *redis.Client
}

func NewGraphQLServer(accountUrl string, catalogURL string, orderURL string, redisURL string) (*Server, error) {
	// Connect to account service
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}

	// Connect to product service
	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	// Connect to order service
	orderClient, err := order.NewClient(orderURL)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	return &Server{
		accountClient,
		catalogClient,
		orderClient,
		redisClient,
	}, nil
}

// Mutation will return MutationResolver
func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}
func (s *Server) Subscription() SubscriptionResolver {
	return &subscriptionResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
