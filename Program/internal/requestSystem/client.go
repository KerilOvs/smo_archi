package requestsystem

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Client represents a client that can submit requests.
type Client struct {
	ID string
}

var requestCounter int
var counterMutex sync.Mutex

// SubmitRequest creates a new request and submits it.
func (c *Client) SubmitRequest(requestType string) *Request {
	counterMutex.Lock()
	requestCounter++
	counter := requestCounter
	counterMutex.Unlock()

	fmt.Printf("Client %s submitted a request of type %s with ID %d\n", c.ID, requestType, counter)
	return &Request{
		ID:        counter,
		Client:    c,
		Status:    "New",
		CreatedAt: time.Now(), // Устанавливаем время создания заявки
	}
}

// GenerateRequests generates requests with a uniform distribution.
func GenerateRequests(client *Client, count int) []*Request {
	requests := make([]*Request, count)
	for i := 0; i < count; i++ {
		requests[i] = client.SubmitRequest("TypeA")
		// Simulate uniform distribution by sleeping for a random time
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
	return requests
}
