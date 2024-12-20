package requestsystem

import (
	"fmt"
	"sort"
	"sync"
)

// RetrievalManager manages the retrieval of requests and assignment to specialists.
type RetrievalManager struct {
	Buffer                 *Buffer
	Specialists            []*Specialist
	CurrentRequest         *Request
	CurrentSpecialistIndex int // Указатель на текущего специалиста в кольцевом буфере
	wg                     sync.WaitGroup
	mu                     sync.Mutex
}

// SelectRequestClick selects a request from the buffer with priority by source number.
func (rm *RetrievalManager) SelectRequestClick() *Request {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Get the next request from the buffer
	request := rm.Buffer.GetNextRequest()
	if request == nil {
		return nil
	}

	// Filter out nil elements from the buffer
	var filteredRequests []*Request
	for _, r := range rm.Buffer.Requests {
		if r != nil {
			filteredRequests = append(filteredRequests, r)
		}
	}

	// Sort requests by source priority (assuming lower ID is higher priority)
	sort.Slice(filteredRequests, func(i, j int) bool {
		return filteredRequests[i].Client.ID < filteredRequests[j].Client.ID
	})

	return request
}

// SendRequestForProcessing sends a request to a specialist for processing and returns the processing time.
func (rm *RetrievalManager) SendRequestForProcessing(request *Request, specialist *Specialist) {
	rm.wg.Add(1) // Increment the WaitGroup counter

	go func() {
		defer rm.wg.Done() // Decrement the WaitGroup counter when done

		// Remove the request from the buffer only after the specialist starts processing it
		rm.Buffer.RemoveRequest(request)

		specialist.TakeRequest(request)
		specialist.ProcessRequest()

	}()
}

// SelectAvailableSpecialist selects an available specialist in a round-robin fashion.
func (rm *RetrievalManager) SelectAvailableSpecialist() *Specialist {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Проходим по кольцу специалистов, начиная с текущего указателя
	for i := 0; i < len(rm.Specialists); i++ {
		specialist := rm.Specialists[rm.CurrentSpecialistIndex]
		rm.CurrentSpecialistIndex = (rm.CurrentSpecialistIndex + 1) % len(rm.Specialists)
		if specialist.IsAvailable() {
			return specialist
		}
	}
	return nil
}

// WaitForSpecialist waits for an available specialist.
func (rm *RetrievalManager) WaitForSpecialist() {
	fmt.Print("")
}

// CheckSpecialistAvailability checks if any specialist is available.
func (rm *RetrievalManager) CheckSpecialistAvailability() bool {
	for _, s := range rm.Specialists {
		if s.IsAvailable() {
			return true
		}
	}
	return false
}

// WaitForAllRequests waits for all requests to be processed.
func (rm *RetrievalManager) WaitForAllRequests() {
	rm.wg.Wait()
}

// PrintSpecialists prints the list of specialists and their current status.
func (rm *RetrievalManager) PrintSpecialists() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	fmt.Println("[ Specialists List: ]")
	for i, specialist := range rm.Specialists {
		status := "Available"
		if !specialist.IsAvailable() {
			status = fmt.Sprintf("Busy with Request ID: %d", specialist.CurrentRequest.ID)
		}
		fmt.Printf("Specialist %d: %s || ", i+1, status)
	}
	fmt.Println("")
}
