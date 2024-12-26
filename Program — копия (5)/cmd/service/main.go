package main

import (
	"fmt"
	"os"
	requestsystem "program/internal/requestSystem"
	"strconv"
	"sync"
	"time"
)

func main() {
	lamb := 200.0    //8.5     // ms время равномерной генерации заявок. от 5 до 9 мс
	lamb_ex := 1.901 // коэфф для експ распределения времени работы. 3 - норм, 5 - быстро
	lamb_ex2 := 1.005

	duration1 := 30 * time.Second // Устанавливаем время работы генератора и процессора
	duration2 := 60 * time.Second

	pause_gen := 100 * time.Millisecond
	pause_proc := 100 * time.Millisecond

	clients_num := 20
	specs_num1 := 2
	specs_num2 := 1
	buffer_cap := 10

	// Создаем файл для вывода в консоль
	consoleLogFile, err := os.Create("console.log")
	if err != nil {
		fmt.Println("Error creating console log file:", err)
		return
	}
	defer consoleLogFile.Close()

	// Перенаправляем стандартный вывод в файл
	oldStdout := os.Stdout
	os.Stdout = consoleLogFile

	// Создаем несколько клиентов
	clients := []*requestsystem.Client{}

	for i := 1; i <= clients_num; i++ {
		client := &requestsystem.Client{ID: strconv.Itoa(i)}
		clients = append(clients, client)
	}

	// Создаем буфер с емкостью 10
	buffer := requestsystem.NewBuffer(buffer_cap)

	specialists := []*requestsystem.Specialist{}
	createdAtTimes := []time.Time{}

	// Создаем специалистов в цикле и добавляем их в список
	for i := 1; i <= specs_num1; i++ {
		specialist := &requestsystem.Specialist{Available: true, Id: i, Lambda: lamb_ex, CreatedAt: time.Now()}
		specialists = append(specialists, specialist)
		createdAtTimes = append(createdAtTimes, specialist.CreatedAt)
	}
	for i := specs_num1 + 1; i <= specs_num1+specs_num2; i++ {
		specialist := &requestsystem.Specialist{Available: true, Id: i, Lambda: lamb_ex2, CreatedAt: time.Now()}
		specialists = append(specialists, specialist)
		createdAtTimes = append(createdAtTimes, specialist.CreatedAt)
	}

	retrievalManager := &requestsystem.RetrievalManager{
		Buffer:      buffer,
		Specialists: specialists,
	}

	stagingManager := &requestsystem.StagingManager{Buffer: buffer}

	// Создаем StatsManager
	statsManager, err := requestsystem.NewStatsManager("stats1.log", specs_num1+specs_num2)
	if err != nil {
		fmt.Println("Error creating stats manager:", err)
		return
	}
	defer statsManager.Close()

	reportManager := requestsystem.NewReportManager(statsManager)

	// Создаем WaitGroup для ожидания завершения горутин
	var wg sync.WaitGroup

	// Запускаем горутину для генерации заявок
	wg.Add(1)
	requestsystem.StartRequestGeneration(clients, stagingManager, retrievalManager, &wg, lamb, statsManager, duration1, pause_gen)

	// Запускаем горутину для обработки заявок
	wg.Add(1)
	requestsystem.StartRequestProcessing(retrievalManager, &wg, statsManager, duration2, pause_proc)

	// Логируем статистику каждые 100 мс
	go func() {
		for {
			time.Sleep(10 * time.Millisecond)
			statsManager.RecordWorkTime(10 * time.Millisecond)
			statsManager.LogStatistics(len(retrievalManager.Specialists), createdAtTimes)
		}
	}()

	// Ожидаем завершения всех горутин
	wg.Wait()

	// Логируем статистику после завершения работы
	statsManager.LogStatistics(len(retrievalManager.Specialists), createdAtTimes)

	// Генерируем отчеты
	reportManager.GenerateSpecialistReport(specialists, createdAtTimes)
	reportManager.GenerateSystemReport()

	// Возвращаем стандартный вывод в консоль
	os.Stdout = oldStdout
}
