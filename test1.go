package main

import (
 "fmt"
 "sync"
)

type Worker struct {
 ID      int
 channel chan string
}

func (w *Worker) Work(wg *sync.WaitGroup) {
 defer wg.Done()
 for {
  data, ok := <-w.channel
  if !ok {
   fmt.Printf("Worker %d: Channel closed\n", w.ID)
   return
  }
  fmt.Printf("Worker %d: %s\n", w.ID, data)
 }
}

func main() {
 inputChannel := make(chan string)
 workers := make(map[int]*Worker)
 var wg sync.WaitGroup

 // Добавление воркеров
 for i := 0; i < 3; i++ {
  worker := &Worker{
   ID:      i,
   channel: make(chan string),
  }
  workers[i] = worker
  wg.Add(1)
  go worker.Work(&wg)
 }

 // Динамическое добавление воркера
 go func() {
  for i := 3; i < 5; i++ {
   worker := &Worker{
    ID:      i,
    channel: make(chan string),
   }
   workers[i] = worker
   wg.Add(1)
   go worker.Work(&wg)
  }
 }()

 // Отправка данных в канал
 go func() {
  for i := 0; i < 10; i++ {
   inputChannel <- fmt.Sprintf("Data %d", i)
  }
  close(inputChannel)
 }()

 // Распределение данных по воркерам
 for data := range inputChannel {
  for _, worker := range workers {
   worker.channel <- data
  }
 }

 // Закрытие каналов воркеров
 for _, worker := range workers {
  close(worker.channel)
 }

 wg.Wait()
}