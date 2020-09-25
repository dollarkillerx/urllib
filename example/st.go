package main

import (
	//"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"runtime"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	//runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		for {
			time.Sleep(time.Second)
			log.Println(runtime.NumGoroutine())
		}
	}()

	limit := make(chan bool, 100)
	for {
		limit <- true
		//time.Sleep(time.Millisecond * 50)
		go func(limit chan bool) {
			defer func() {
				//log.Println("Over")
				<-limit
			}()

			_, err := http.Get("http://127.0.0.1:8986/test")
			if err != nil {
				log.Println(err)
			}
		}(limit)
	}
}
