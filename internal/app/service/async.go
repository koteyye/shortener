package service

import (
	"strings"
	"sync"

	"github.com/koteyye/shortener/internal/app/models"
)

func add(doneCh chan struct{}, urls []string) chan string {
	addURL := make(chan string)

	go func() {
		defer close(addURL)

		for _, data := range urls {
			select {
			case <-doneCh:
				return
			case addURL <- data:
			}
		}
	}()

	return addURL
}

func validateUser(doneCh chan struct{}, inputCh chan string, urlListByUser []*models.AllURLs) chan string {
	validateURL := make(chan string)

	go func() {
		defer close(validateURL)

		for data := range inputCh {
			var result string
			for _, urlItemByUser := range urlListByUser {
				if strings.Contains(urlItemByUser.ShortURL, data) {
					result = data
					break
				}
			}

			select {
			case <-doneCh:
				return
			case validateURL <- result:
			}
		}
	}()
	return validateURL
}

func fanOut(doneCh chan struct{}, inputCh chan string, urlListByUser []*models.AllURLs) []chan string {
	//количество горутин
	numWorkers := 10

	channels := make([]chan string, numWorkers)

	for i := 0; i < numWorkers; i++ {
		// канал из горутины add
		addResultCh := validateUser(doneCh, inputCh, urlListByUser)
		// отправляем в слайс каналов
		channels[i] = addResultCh
	}

	return channels
}

func fanIn(doneCh chan struct{}, resultChs ...chan string) chan string {
	finalCh := make(chan string)

	var wg sync.WaitGroup

	//перебираем все входящие каналы
	for _, ch := range resultChs {
		chClosure := ch

		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range chClosure {
				select {
				case <-doneCh:
					return
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		//ждем завершения всех горутин
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}
