package requestsystem

import "fmt"

// StagingManager manages the staging of requests.
type StagingManager struct {
	Buffer         *Buffer
	CurrentRequest *Request
}

// InitiatePlacement initiates the placement of a request in the buffer.
func (sm *StagingManager) InitiatePlacement(request *Request) {
	fmt.Println("Initiating placement of request")
}

// CheckIsBufferFull checks if the buffer is full.
func (sm *StagingManager) CheckIsBufferFull() bool {
	return sm.Buffer.IsFull()
}

// AddRequestBuffer adds a request to the buffer.
func (sm *StagingManager) AddRequestBuffer(request *Request) {
	if sm.Buffer.AddRequest(request) {
		fmt.Printf("Request %d added to buffer\n", request.ID)
	} else {
		fmt.Printf("Buffer is full, discarding the last request and add %d \n", request.ID)
	}
}

// RemoveOldest removes the oldest request from the buffer.
func (sm *StagingManager) RemoveOldest() {
	oldest := sm.Buffer.GetNextRequest()
	if oldest != nil {
		fmt.Printf("Removed oldest request %d from buffer\n", oldest.ID)
	}
}
