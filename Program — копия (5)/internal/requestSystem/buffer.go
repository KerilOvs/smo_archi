package requestsystem

import (
	"fmt"
	"math"
	"sync"
)

// Buffer represents a circular buffer for requests.
type Buffer struct {
	Requests []*Request
	Capacity int
	Head     int // Points to the next position to write
	Tail     int // Points to the next position to read
	Full     bool
	mu       sync.Mutex
}

// AddRequest adds a request to the circular buffer.
func (b *Buffer) AddRequest(request *Request) bool {
	b.mu.Lock()
	// defer b.mu.Unlock()

	status := true

	if b.Full {
		// If buffer is full, overwrite the last request (Head - 1)
		b.Head = (b.Head - 1 + b.Capacity) % b.Capacity
		status = false
	}

	b.Requests[b.Head] = request
	b.Head = (b.Head + 1) % b.Capacity

	if b.Head == b.Tail {
		b.Full = true
	}
	b.mu.Unlock()
	return status
}

// RemoveRequest removes a specific request from the buffer.
func (b *Buffer) RemoveRequest(request *Request) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i := b.Tail; i != b.Head; i = (i + 1) % b.Capacity {
		if b.Requests[i] == request {
			// Shift all subsequent requests to the left
			for j := i; j != b.Head; j = (j + 1) % b.Capacity {
				next := (j + 1) % b.Capacity
				b.Requests[j] = b.Requests[next]
			}
			b.Head = (b.Head - 1 + b.Capacity) % b.Capacity
			break
		}
	}
}

// GetNextRequest gets the next request from the circular buffer and removes it.
func (b *Buffer) GetNextRequest() *Request {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.Full && b.Head == b.Tail {
		return nil // Buffer is empty
	}

	request := b.Requests[b.Tail]
	b.Requests[b.Tail] = nil // Удаляем заявку из буфера
	b.Tail = (b.Tail + 1) % b.Capacity
	b.Full = false

	return request
}

// IsFull checks if the circular buffer is full.
func (b *Buffer) IsFull() bool {
	return b.Full
}

func (b *Buffer) IsEmpty() bool {
	// Если указатель Head равен указателю Tail и буфер не полон, значит буфер пуст
	return b.Head == b.Tail && !b.Full
}

// PrintBufferContent prints the current content of the buffer.
func (b *Buffer) PrintBufferContent() {
	b.mu.Lock()
	defer b.mu.Unlock()

	fmt.Println("[ Buffer Content: ]", "{ ", math.Abs(float64(b.Head-b.Tail)), "}")
	if b.Head == b.Tail && !b.Full {
		fmt.Println("Buffer is empty")
		return
	}

	// Если буфер полон, выводим содержимое
	if b.Full {
		fmt.Println("Buffer is full")
	}

	// Выводим содержимое буфера, начиная с Tail и заканчивая Head
	// Для циклического буфера нужно учитывать, что Head может быть меньше Tail
	for i := 0; i < b.Capacity; i++ {
		index := (b.Tail + i) % b.Capacity
		if b.Requests[index] != nil {
			fmt.Printf("Request ID: %d, Client ID: %s  || ", b.Requests[index].ID, b.Requests[index].Client.ID)
		}
	}
	fmt.Println("")
}

// NewBuffer creates a new circular buffer with the given capacity.
func NewBuffer(capacity int) *Buffer {
	return &Buffer{
		Requests: make([]*Request, capacity),
		Capacity: capacity,
		Head:     0,
		Tail:     0,
		Full:     false,
	}
}
