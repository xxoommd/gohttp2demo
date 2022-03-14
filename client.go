package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"

	"golang.org/x/net/http2"
)

const url = "https://localhost:8001"

var wg sync.WaitGroup

func main() {
	num := 100
	for i := 0; i < 1000; i++ {
		wg.Add(num)
		go testHTTP(num)
	}

	wg.Wait()
}

func testHTTP(n int) {
	for i := 0; i < n; i++ {
		client := &http.Client{}
		client.Transport = &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			}}

		resp, err := client.Get(url)
		if err != nil {
			log.Fatal("get url err:", err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed reading response body: %s", err)
		}

		fmt.Printf("Got response %d: %s %s\n", resp.StatusCode, resp.Proto, string(body))
		wg.Done()
	}
}
