package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func main() {
	wg := &sync.WaitGroup{}

	for i := 0; i < 50; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go get((i+1)*(j+1)-1, wg)
		}

		time.Sleep(time.Millisecond * 500)
	}

	wg.Wait()
}

func get(i int, wg *sync.WaitGroup) {
	proxyUrl, _ := url.Parse("http://localhost:32199")

	client := http.Client{
		Timeout: time.Minute * 30,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	res, err := client.Get("https://example.com/")

	if res != nil {
		fmt.Printf("%d: %d\n", i, res.StatusCode)
	}

	if err != nil {
		fmt.Printf("%d: %s\n", i, err)
		wg.Done()
		return
	}

	if res.StatusCode != http.StatusOK {
		bodyData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("%d: %s\n", i, err)
			wg.Done()
			return
		}

		fmt.Println(string(bodyData))
	}

	err = res.Body.Close()
	if err != nil {
		fmt.Printf("%d: %s\n", i, err)
		wg.Done()
		return
	}

	wg.Done()
}
