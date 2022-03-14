package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/http2"
)

const url = "https://localhost:8001/"

// var wg sync.WaitGroup

func main() {
	// num := 2
	// for i := 0; i < 2; i++ {
	// 	wg.Add(num)
	// 	go testHTTP(num)
	// }

	// wg.Wait()
	start := time.Now()
	testHTTP(100)
	fmt.Println("Time:", time.Since(start))
}

func testHTTP(n int) {
	wg2 := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		client := &http.Client{}
		client.Transport = &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			}}

		for j := 0; j < 100; j++ {
			wg2.Add(1)
			indexI, indexJ := i, j
			go func() {
				fullURL := fmt.Sprintf("%s?i=%d&j=%d", url, indexI, indexJ)
				resp, err := client.Get(fullURL)
				if err != nil {
					fmt.Println("get url err:", err)
				} else {
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Printf("Failed reading response body: %s\n", err)
					} else {
						fmt.Printf("[%d-%d]Got response %d: %s %s\n", indexI, indexJ, resp.StatusCode, resp.Proto, string(body))
					}
				}
				resp.Body.Close()
				wg2.Done()
			}()
		}
	}
	wg2.Wait()
}
