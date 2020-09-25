package main

import (
	"github.com/dollarkillerx/urllib"
	"log"
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

	limit := make(chan bool, 300)
	for {
		limit <- true
		//time.Sleep(time.Millisecond * 50)
		go func(limit chan bool) {
			defer func() {
				//log.Println("Over")
				<-limit
			}()

			_, _, err := urllib.Get("http://127.0.0.1:8986/test").Byte()
			if err != nil {
				log.Println(err)
				return
			}

			//_, err := urllib.DefGet("http://127.0.0.1:8986/test")
			//if err != nil {
			//	log.Println(err)
			//	return
			//}

			//_, _, err := fasthttp.Get(nil, "http://127.0.0.1:8986/test")
			//if err != nil {
			//	log.Println(err)
			//	return
			//}

			//req := fasthttp.AcquireRequest()
			//resp := fasthttp.AcquireResponse()
			//defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
			//defer fasthttp.ReleaseResponse(resp) // <- do not forget to release
			//
			//req.Header.SetMethod("GET")
			//req.SetRequestURI("http://127.0.0.1:8986/test")
			//
			//client := fasthttp.Client{}
			//
			//err := client.Do(req, resp)
			//if err != nil {
			//	log.Println(err)
			//}
			//fasthttp.ReleaseRequest(req)
			//fasthttp.ReleaseResponse(resp)

			//_, _, err := fasthttp.GetTimeout(nil, "http://127.0.0.1:8986/test", time.Second)
			//if err != nil {
			//	log.Println(err)
			//}

		}(limit)
	}
}
