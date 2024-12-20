package requestsystem

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var outputMutex sync.Mutex // Мьютекс для синхронизации вывода

// Specialist represents a specialist that can process requests.
type Specialist struct {
	CurrentRequest *Request
	Available      bool
	WorkTime       time.Duration
	Lambda         float64 // 3 - более-менее, чем меньше, тем быстрее
	Id             int
	mu             sync.Mutex
}

// TakeRequest assigns a request to the specialist.
func (s *Specialist) TakeRequest(request *Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CurrentRequest = request
	s.Available = false
}

// ProcessRequest processes the current request.
func (s *Specialist) ProcessRequest() {
	outputMutex.Lock()
	fmt.Printf("Specialist %d Processing request %d\n", s.Id, s.CurrentRequest.ID)
	outputMutex.Unlock()
	s.CurrentRequest.UpdateStatus("Processing")

	// Simulate exponential distribution for processing time
	processingTime := time.Duration(rand.ExpFloat64() / s.Lambda * float64(time.Second))
	s.WorkTime = processingTime
	time.Sleep(processingTime)

	outputMutex.Lock()
	fmt.Printf("Request %d completed by spec %d\n", s.CurrentRequest.ID, s.Id)
	outputMutex.Unlock()
	s.CurrentRequest.UpdateStatus("Completed")
	s.CurrentRequest = nil
	s.Available = true

}

// IsAvailable checks if the specialist is available.
func (s *Specialist) IsAvailable() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Available
}
