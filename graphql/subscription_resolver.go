package main

import (
	"context"
	"sync"
)

type subscriptionResolver struct {
	server          *Server
	messageChannels map[string]chan *string
	userChannels    map[string]chan string
	mutex           sync.Mutex
}

func (r *subscriptionResolver) MessagePosted(ctx context.Context, id string, message *string) (<-chan string, error) {
	c := make(chan string)
	c <- "done"
	return c, nil

}

func (r *subscriptionResolver) UserJoined(ctx context.Context, id string) (<-chan string, error) {
	strID := id
	err := r.createUser(strID)
	if err != nil {
		return nil, err
	}

	// Create new channel for request
	users := make(chan string, 1)
	r.mutex.Lock()
	r.userChannels[strID] = users
	r.mutex.Unlock()

	// Delete channel when done
	go func() {
		<-ctx.Done()
		r.mutex.Lock()
		delete(r.userChannels, strID)
		r.mutex.Unlock()
	}()

	return users, nil
}

func (r *subscriptionResolver) createUser(userID string) error {
	// Upsert user
	if err := r.server.redisClient.SAdd("users", userID).Err(); err != nil {
		return err
	}
	// Notify new user joined
	r.mutex.Lock()
	for _, ch := range r.userChannels {
		ch <- userID
	}
	r.mutex.Unlock()
	return nil
}
