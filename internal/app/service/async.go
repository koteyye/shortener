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

func validateUser(doneCh chan struct{}, urlListByUser []*models.URLList, urls []string) chan string {
	validateURL := make(chan string)

	go func() {
		defer close(validateURL)

		for _, url := range urls {
			var result string
			for _, urlItem := range urlListByUser {
				if strings.Contains(urlItem.ShortURL, url) {
					result = url
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
