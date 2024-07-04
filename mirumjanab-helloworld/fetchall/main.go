package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Лол")
		return
	}
	var wg sync.WaitGroup

	var duration time.Duration
	for _, url := range os.Args[1:] {

		wg.Add(1)

		go func(url string) {
			defer wg.Done()
			start := time.Now()
			res, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				return
			}

			if res.StatusCode == http.StatusOK {
				fmt.Println("The requested URL was found")
			} else if res.StatusCode == http.StatusNotFound {
				fmt.Println("The requested URL was not found")
			}

			body, err := io.ReadAll(res.Body)

			if err != nil {
				fmt.Println(err)
				return
			}

			size := len(body)

			res.Body.Close()
			end := time.Now()
			duration = end.Sub(start)
			fmt.Println(duration, size, url)

		}(url)
	}
	wg.Wait()
	fmt.Println(duration)
}
