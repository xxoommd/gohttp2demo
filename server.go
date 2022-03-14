package main

import (
	"log"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// var wg sync.WaitGroup

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		i := r.Form.Get("i")
		j := r.Form.Get("j")
		log.Printf("Got connection [%s-%s]: %s", i, j, r.Proto)
		w.Write([]byte("Hello"))
	})

	srv := &http.Server{
		Addr:    ":8001",
		Handler: h2c.NewHandler(handler, &http2.Server{}),
	}

	log.Printf("Serving on https://%s", srv.Addr)
	// wg.Add(1)
	// go func() {

	// }()

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("start http server fail: ", err)
	}
	// wg.Wait()
}
