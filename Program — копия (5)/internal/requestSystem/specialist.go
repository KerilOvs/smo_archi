package requestsystem

import (
	"fmt"
	"math"
	"sync"
	"time"
)

var outputMutex sync.Mutex // Мьютекс для синхронизации вывода

// Specialist represents a specialist that can process requests.
type Specialist struct {
	CurrentRequest         *Request
	Available              bool
	WorkTime               time.Duration
	Lambda                 float64
	ProcessedRequestsCount int // Количество отработанных заявок
	Id                     int
	CreatedAt              time.Time
	mu                     sync.Mutex
}

// TakeRequest assigns a request to the specialist.
func (s *Specialist) TakeRequest(request *Request) {
	// s.mu.Lock()
	// defer s.mu.Unlock() // поробовать перевести в конец, мб пропадет хуйня
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
	processingTime := time.Duration(10 * math.Exp(s.Lambda*float64(s.ProcessedRequestsCount)) * float64(time.Millisecond))
	s.WorkTime = time.Duration(float64(processingTime) * 2.7)
	time.Sleep(processingTime)

	outputMutex.Lock()
	if s.CurrentRequest != nil {
		fmt.Printf("Request %d completed by spec %d\n", s.CurrentRequest.ID, s.Id)
		s.CurrentRequest.UpdateStatus("Completed")
		s.CurrentRequest = nil
	} else {
		fmt.Println("Request-huinya")
	}
	s.Available = true
	s.ProcessedRequestsCount++ // Увеличиваем счетчик обработанных заявок

	outputMutex.Unlock()
}

// IsAvailable checks if the specialist is available.
func (s *Specialist) IsAvailable() bool {
	// s.mu.Lock()
	// defer s.mu.Unlock()
	return s.Available
}
