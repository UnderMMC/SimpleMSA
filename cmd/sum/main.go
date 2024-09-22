package main

import (
	"github.com/SimpleMSA/internal/app"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2) // Указываем, что будет два сервиса

	go app.New().Run(&wg)      // Запускаем первый сервис в горутине
	go app.NewOrder().Run(&wg) // Запускаем второй сервис в горутине

	wg.Wait() // Ожидаем завершения обоих сервисов
	//app.New().Run()
	//app.NewOrder().Run()
}
