package main

import (
	"fmt"
	requestsystem "program/internal/requestSystem"
	"strconv"
	"sync"
	"time"
)

func main() {
	lamb := 0.01   // 10 ms время равномерной генерации заявок
	lamb_ex := 5.0 // коэфф для експ распределения времени работы. 3 - норм, 5 - быстро

	duration1 := 10 * time.Second // Устанавливаем время работы генератора и процессора
	duration2 := 30 * time.Second

	clients_num := 3
	specs_num := 10
	buffer_cap := 20

	// Создаем несколько клиентов
	clients := []*requestsystem.Client{
		{ID: "1"},
		{ID: "2"},
		{ID: "3"},
	}

	for i := 1; i <= clients_num; i++ {
		client := &requestsystem.Client{ID: strconv.Itoa(i)}
		clients = append(clients, client)
	}

	// Создаем буфер с емкостью 10
	buffer := requestsystem.NewBuffer(buffer_cap)

	specialists := []*requestsystem.Specialist{}

	// Создаем специалистов в цикле и добавляем их в список
	for i := 1; i <= specs_num; i++ {
		specialist := &requestsystem.Specialist{Available: true, Id: i, Lambda: lamb_ex}
		specialists = append(specialists, specialist)
	}

	retrievalManager := &requestsystem.RetrievalManager{
		Buffer:      buffer,
		Specialists: specialists,
	}

	stagingManager := &requestsystem.StagingManager{Buffer: buffer}

	// Создаем StatsManager
	statsManager, err := requestsystem.NewStatsManager("stats1.log", specs_num)
	if err != nil {
		fmt.Println("Error creating stats manager:", err)
		return
	}
	defer statsManager.Close()

	// Создаем WaitGroup для ожидания завершения горутин
	var wg sync.WaitGroup

	// Запускаем горутину для генерации заявок
	wg.Add(1)
	requestsystem.StartRequestGeneration(clients, stagingManager, retrievalManager, &wg, lamb, statsManager, duration1)

	// Запускаем горутину для обработки заявок
	wg.Add(1)
	requestsystem.StartRequestProcessing(retrievalManager, &wg, statsManager, duration2)

	// Логируем статистику каждые 100 мс
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			statsManager.LogStatistics(len(retrievalManager.Specialists))
		}
	}()

	// Ожидаем завершения всех горутин
	wg.Wait()

	// Логируем статистику после завершения работы
	statsManager.LogStatistics(len(retrievalManager.Specialists))
}
