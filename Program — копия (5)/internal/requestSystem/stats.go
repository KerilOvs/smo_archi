package requestsystem

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// StatsManager manages the collection and logging of statistics.
type StatsManager struct {
	TotalRequests       int
	RejectedRequests    int
	TotalBufferTime     time.Duration
	TotalProcessingTime time.Duration
	SpecialistUsage     map[int]int // Map of specialist ID to the number of requests they processed
	mu                  sync.Mutex
	File                *os.File
	LastLogTime         time.Time
	logChannel          chan string           // Буферизованный канал для записи логов
	TotalSystemTime     time.Duration         // Общее время работы системы
	SpecialistWorkTime  map[int]time.Duration // Время работы каждого специалиста
}

// NewStatsManager creates a new StatsManager and initializes the log file.
func NewStatsManager(filename string, spec_num int) (*StatsManager, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	// Write CSV header
	str := "Timestamp,TotalRequests,RejectedRequests,ProbabilityOfRejection,AverageBufferTime,AverageProcessingTime"
	// for i := 0; i < spec_num; i++ {
	// 	str += fmt.Sprintf(",Specialist%dLoad", i+1)
	// }

	for i := 0; i < spec_num; i++ {
		str += fmt.Sprintf(",Specialist%dWorkTimeRatio", i+1)
	}
	str += "\n"

	_, err = file.WriteString(str)
	if err != nil {
		return nil, err
	}

	sm := &StatsManager{
		SpecialistUsage:    make(map[int]int),
		SpecialistWorkTime: make(map[int]time.Duration),
		File:               file,
		LastLogTime:        time.Now(),
		logChannel:         make(chan string, 100), // Буферизованный канал
	}

	// Запуск горутины для записи логов в файл
	go sm.logWriter()

	return sm, nil
}

// RecordRequest records a new request and updates the total request count.
func (sm *StatsManager) RecordRequest() {
	//sm.mu.Lock()
	// defer sm.mu.Unlock()
	sm.TotalRequests++
	//sm.mu.Unlock()
}

// RecordRejectedRequest records a rejected request.
func (sm *StatsManager) RecordRejectedRequest() {
	//sm.mu.Lock()
	// defer sm.mu.Unlock()
	sm.RejectedRequests++
	//sm.mu.Unlock()
}

// RecordBufferTime records the time a request spent in the buffer.
func (sm *StatsManager) RecordBufferTime(duration time.Duration) {
	sm.mu.Lock()
	// defer sm.mu.Unlock()
	sm.TotalBufferTime += duration
	sm.mu.Unlock()
}

// RecordProcessingTime records the time a request spent being processed.
func (sm *StatsManager) RecordProcessingTime(duration time.Duration) {
	sm.mu.Lock()
	// defer sm.mu.Unlock()
	sm.TotalProcessingTime += duration
	sm.mu.Unlock()
}

// RecordSpecialistUsage records the usage of a specialist.
func (sm *StatsManager) RecordSpecialistUsage(specialistID int) {
	sm.mu.Lock()
	// defer sm.mu.Unlock()
	sm.SpecialistUsage[specialistID]++
	sm.mu.Unlock()
}

func (sm *StatsManager) RecordSpecialistWorkTime(specialistID int, workTime time.Duration) {
	sm.mu.Lock()
	// defer sm.mu.Unlock()
	sm.SpecialistWorkTime[specialistID] += workTime
	// sm.TotalSystemTime += workTime
	sm.mu.Unlock()
}

func (sm *StatsManager) RecordWorkTime(workTime time.Duration) {
	sm.mu.Lock()
	// defer sm.mu.Unlock()
	sm.TotalSystemTime += workTime
	sm.mu.Unlock()
}

// CalculateProbabilityOfRejection calculates the probability of rejection.
func (sm *StatsManager) CalculateProbabilityOfRejection() float64 {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.TotalRequests == 0 {
		return 0.0
	}
	return float64(sm.RejectedRequests) / float64(sm.TotalRequests)
}

// CalculateAverageBufferTime calculates the average time a request spends in the buffer.
func (sm *StatsManager) CalculateAverageBufferTime() float64 {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.TotalRequests == 0 {
		return 0.0
	}
	return float64(sm.TotalBufferTime.Nanoseconds()) / float64(sm.TotalRequests-sm.RejectedRequests) / 1e6 // Convert to milliseconds
}

// CalculateAverageProcessingTime calculates the average time a request spends being processed.
func (sm *StatsManager) CalculateAverageProcessingTime() float64 {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.TotalRequests == 0 {
		return 0.0
	}
	return float64(sm.TotalProcessingTime.Nanoseconds()) / float64(sm.TotalRequests-sm.RejectedRequests) / 1e6 // Convert to milliseconds
}

// CalculateSpecialistLoad calculates the load of each specialist.
func (sm *StatsManager) CalculateSpecialistLoad(totalSpecialists int) map[int]float64 {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	load := make(map[int]float64)
	for id, count := range sm.SpecialistUsage {
		load[id] = float64(count) / float64(sm.TotalRequests-sm.RejectedRequests)
	}
	return load
}

// LogStatistics logs the current statistics to the file.
func (sm *StatsManager) LogStatistics(totalSpecialists int, createdAtTimes []time.Time) {
	// Check if 100ms have passed since the last log
	if time.Since(sm.LastLogTime) < 100*time.Millisecond {
		return
	}

	probRejection := sm.CalculateProbabilityOfRejection()
	avgBufferTime := sm.CalculateAverageBufferTime()
	avgProcessingTime := sm.CalculateAverageProcessingTime()
	// specialistLoad := sm.CalculateSpecialistLoad(totalSpecialists)

	// Prepare the log entry
	logEntry := fmt.Sprintf("%s,%d,%d,%.4f,%.6f,%.6f",
		time.Now().Format(time.RFC3339Nano),
		sm.TotalRequests,
		sm.RejectedRequests,
		probRejection,
		avgBufferTime,
		avgProcessingTime,
	)

	// Add specialist loads
	// for i := 1; i <= totalSpecialists; i++ {
	// 	logEntry += fmt.Sprintf(",%.4f", specialistLoad[i])
	// }

	// Add specialist loads
	for i := 1; i <= totalSpecialists; i++ {
		specialistWorkTimeRatio := float64(sm.SpecialistWorkTime[i]) / float64(time.Since(createdAtTimes[i-1]))
		if specialistWorkTimeRatio >= 1.0 {
			logEntry += fmt.Sprintf(",%.4f", 1.0)
		} else {
			logEntry += fmt.Sprintf(",%.4f", specialistWorkTimeRatio)
		}

	}

	logEntry += "\n"

	// Send log entry to the channel
	sm.logChannel <- logEntry

	// Update the last log time
	sm.LastLogTime = time.Now()
}

// logWriter writes log entries from the channel to the file.
func (sm *StatsManager) logWriter() {
	for entry := range sm.logChannel {
		_, err := sm.File.WriteString(entry)
		if err != nil {
			fmt.Println("Error writing to log file:", err)
		}
	}
}

// Close closes the log file.
func (sm *StatsManager) Close() {
	close(sm.logChannel) // Закрываем канал, чтобы завершить горутину logWriter
	sm.File.Close()
}
