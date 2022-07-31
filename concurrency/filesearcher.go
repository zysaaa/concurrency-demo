package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
)

const tofind = "select-test.go"

func main() {
	wg := sync.WaitGroup{}
	search("/Users/user/Documents", &wg)
	wg.Add(1)
	wg.Wait()
}

var permit = make(chan struct{}, 30)

func search(dir string, wg *sync.WaitGroup) {
	permit <- struct{}{}

	defer func() {
		<-permit
	}()
	defer func() {
		wg.Done()
	}()

	items, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, item := range items {
		if item.IsDir() {
			wg.Add(1)
			go search(filepath.Join(dir, item.Name()), wg)
		} else {
			if item.Name() == tofind {
				fmt.Println(filepath.Join(dir, item.Name()))
			}
		}
	}
}
