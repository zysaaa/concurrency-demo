package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const jobCount = 8
	jobs := make(chan int, jobCount)
	results := make(chan int, jobCount)

	for w := 1; w <= 3; w++ {
		go work(w, jobs, results)
	}
	for j := 1; j <= jobCount; j++ {
		jobs <- j
	}
	close(jobs)
	fmt.Println("Closed job")
	for a := 1; a <= jobCount; a++ {
		<-results
	}
	close(results)
}

func work(id int, jobs <-chan int, results chan<- int) {
	var wg sync.WaitGroup
	for j := range jobs {
		wg.Add(1)
		go func(job int) {
			fmt.Println("worker", id, "started job", job)
			time.Sleep(time.Second)
			fmt.Println("worker", id, "finished job", job)
			results <- job * 2
			wg.Done()
		}(j)
	}
	wg.Wait()
}
