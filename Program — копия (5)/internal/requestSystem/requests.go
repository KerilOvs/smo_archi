package requestsystem

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// StartRequestGeneration запускает горутину для генерации заявок с ограничением по времени
func StartRequestGeneration(clients []*Client, stagingManager *StagingManager, retrievalManager *RetrievalManager, wg *sync.WaitGroup, lamb float64, statsManager *StatsManager, duration time.Duration, pauseTime time.Duration) {
	go func() {
		defer wg.Done()
		timer := time.NewTimer(duration)
		defer timer.Stop()

		for {
			select {
			case <-timer.C:
				// Таймер истек, завершаем горутину
				return
			default:
				// Случайно выбираем клиента
				client := clients[rand.Intn(len(clients))]

				// Создаем заявку
				request := client.SubmitRequest("TypeA")
				// Записываем статистику о новой заявке
				statsManager.RecordRequest()

				// Добавляем заявку в буфер
				if retrievalManager.CheckSpecialistAvailability() {
					availableSpecialist := retrievalManager.SelectAvailableSpecialist()
					if availableSpecialist != nil {
						retrievalManager.SendRequestForProcessing(request, availableSpecialist, false)
					} else {
						fmt.Println("pizdec")
					}
				} else {
					if !stagingManager.Buffer.AddRequest(request) {
						// Если буфер полон, записываем отклоненную заявку
						statsManager.RecordRejectedRequest()
					}
				}

				// Выводим содержимое буфера
				stagingManager.Buffer.PrintBufferContent()
				retrievalManager.PrintSpecialists()

				// Ожидаем случайное время (интенсивность от 0 до 1)
				time.Sleep(time.Duration(lamb * float64(time.Millisecond)))
			}
		}
	}()
}

// StartRequestProcessing запускает горутину для обработки заявок с ограничением по времени
func StartRequestProcessing(retrievalManager *RetrievalManager, wg *sync.WaitGroup, statsManager *StatsManager, duration time.Duration, pauseTime time.Duration) {
	go func() {
		defer wg.Done()
		timer := time.NewTimer(duration)
		defer timer.Stop()

		for {
			select {
			case <-timer.C:
				// Таймер истек, завершаем горутину
				return
			default:
				// Выбираем следующую заявку из буфера
				nextRequest := retrievalManager.SelectRequestClick()
				if nextRequest == nil {
					// Если буфер пуст, ждем
					time.Sleep(pauseTime)
					continue
				}

				// Выбираем доступного специалиста
				availableSpecialist := retrievalManager.SelectAvailableSpecialist()
				if availableSpecialist == nil {
					// Если нет доступных специалистов, ждем
					time.Sleep(pauseTime)
					continue
				}

				// Записываем время, проведенное в буфере
				statsManager.RecordBufferTime(time.Since(nextRequest.CreatedAt))

				// Отправляем заявку на обработку
				retrievalManager.SendRequestForProcessing(nextRequest, availableSpecialist, true)
				statsManager.RecordProcessingTime(availableSpecialist.WorkTime)

				// Записываем использование специалиста
				statsManager.RecordSpecialistUsage(availableSpecialist.Id)
				statsManager.RecordSpecialistWorkTime(availableSpecialist.Id, availableSpecialist.WorkTime)

				// Ожидаем некоторое время перед следующей итерацией
				time.Sleep(pauseTime)
			}
		}
	}()
}
