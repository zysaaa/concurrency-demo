package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	roots := []string{"/"}
	tick := time.Tick(1500 * time.Millisecond)
	wg := sync.WaitGroup{}
	fileSize := make(chan int64)
	for _, root := range roots {
		wg.Add(1)
		go walk(root, &wg, fileSize)
	}
	go func() {
		wg.Wait()
		close(fileSize)
	}()
	var files, bytes int64
loop:
	for {
		select {
		case size, ok := <-fileSize:
			if !ok {
				break loop
			}
			files++
			bytes += size
		case <-tick:
			printDiskUsage(files, bytes)
		}
	}
	printDiskUsage(files, bytes)
}

func printDiskUsage(files, bytes int64) {
	fmt.Printf("%d files  %.1f GB\n", files, float64(bytes)/1e9)
}

var semaphore = make(chan struct{}, 30)

func walk(dir string, wg *sync.WaitGroup, ch chan<- int64) {
	semaphore <- struct{}{}        // acquire token
	defer func() { <-semaphore }() // release token
	defer wg.Done()
	items, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, item := range items {
		if item.IsDir() {
			wg.Add(1)
			go walk(filepath.Join(dir, item.Name()), wg, ch)
		} else {
			ch <- item.Size()
		}
	}
}
