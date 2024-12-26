package requestsystem

import (
	"fmt"
	"time"
)

// ReportManager отвечает за формирование и вывод отчетов
type ReportManager struct {
	StatsManager *StatsManager
}

// NewReportManager создает новый ReportManager
func NewReportManager(statsManager *StatsManager) *ReportManager {
	return &ReportManager{
		StatsManager: statsManager,
	}
}

// GenerateSpecialistReport генерирует отчет по каждому специалисту
func (rm *ReportManager) GenerateSpecialistReport(specialists []*Specialist, createdAtTimes []time.Time) {
	fmt.Println("Stats for Specialists:")
	fmt.Printf("%-5s %-15s %-15s %-20s %-15s %-15s\n", "ID", "WorkTime", "Lambda", "ProcessedRequests", "LoadPercentage", "LoadPercentageByTime")

	for _, specialist := range specialists {
		workTime := specialist.WorkTime
		lambda := specialist.Lambda
		processedRequests := specialist.ProcessedRequestsCount
		loadPercentage := float64(processedRequests) / float64(rm.StatsManager.TotalRequests-rm.StatsManager.RejectedRequests) * 100
		// specialistWorkTimeRatio := float64(sm.SpecialistWorkTime[i]) / float64(time.Since(createdAtTimes[i-1]))
		LoadPercentageByTime := float64(rm.StatsManager.SpecialistWorkTime[specialist.Id-1]) / float64(time.Since(createdAtTimes[specialist.Id-1]))

		fmt.Printf("%-5d %-15s %-15.4f %-20d %-15.2f%-15.2f%%\n", specialist.Id, workTime, lambda, processedRequests, loadPercentage, LoadPercentageByTime)
	}
}

// GenerateSystemReport генерирует отчет по системе
func (rm *ReportManager) GenerateSystemReport() {
	fmt.Println("\nStats for System:")
	fmt.Printf("%-20s %-20s %-20s %-20s %-20s\n", "TotalRequests", "RejectedRequests", "TotalBufferTime", "TotalProcessingTime", "TotalSystemTime")

	totalRequests := rm.StatsManager.TotalRequests
	rejectedRequests := rm.StatsManager.RejectedRequests
	totalBufferTime := rm.StatsManager.TotalBufferTime
	totalProcessingTime := rm.StatsManager.TotalProcessingTime
	totalSystemTime := rm.StatsManager.TotalSystemTime

	fmt.Printf("%-20d %-20d %-20s %-20s %-20s\n", totalRequests, rejectedRequests, totalBufferTime, totalProcessingTime, totalSystemTime)
}
