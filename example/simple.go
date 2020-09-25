package main

import (
	//"github.com/dollarkillerx/urllib"

	"github.com/dollarkillerx/fasthttp"
	"log"
	"net/http"
	"runtime"
	"runtime/pprof"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	p := pprof.Lookup("goroutine")
	p.WriteTo(w, 1)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	//runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		for {
			time.Sleep(time.Second)
			log.Println(runtime.NumGoroutine())
		}
	}()
	go func() {
		http.HandleFunc("/", handler)
		http.ListenAndServe(":11181", nil)
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

			//client := http.Client{
			//	//Transport: &http.Transport{
			//	//	TLSClientConfig: &tls.Config{
			//	//		InsecureSkipVerify: true,
			//	//	},
			//	//},
			//}
			//url, err := url.Parse("http://0.0.0.0:8986/test")
			//if err != nil {
			//	log.Println(err)
			//	return
			//}
			//_, err = client.Do(&http.Request{
			//	URL:        url,
			//	Method:     "GET",
			//	Proto:      "HTTP/1.1",
			//	//ProtoMajor: 1,
			//	//ProtoMinor: 1,
			//})
			//if err != nil {
			//	log.Println(err)
			//	return
			//}

			_, _, err := fasthttp.Get(nil, "https://0.0.0.0:8091/")
			if err != nil {
				log.Println(err)
				return
			}

			//fasthttp.

		}(limit)
	}
}

//func BasicPprofTime(timeout time.Duration) {
//	cpuf, e := os.Create("cpu_profile")
//	if e != nil {
//		log.Fatalln(e)
//	}
//
//	pprof.StartCPUProfile(cpuf)
//
//	defer pprof.StopCPUProfile()
//
//	time.Sleep(timeout)
//
//	log.Println("关闭分析|||关闭分析|||关闭分析|||关闭分析|||关闭分析|||关闭分析|||关闭分析|||关闭分析|||关闭分析|||关闭分析|||  CPU")
//
//	memf, err := os.Create("mem_profile")
//	if err != nil {
//		log.Fatal("could not create memory profile: ", err)
//	}
//	if err := pprof.WriteHeapProfile(memf); err != nil {
//		log.Fatal("could not write memory profile: ", err)
//	}
//	memf.Close()
//
//	log.Println("关闭分析|||关闭分析|||关闭分析|||关闭分析|||关闭分析|||关闭分析|||关闭分析|||关闭分析|||关闭分析|||关闭分析||| MEMF")
//}

func doRequest(url string) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	//req.SetRequestURI(url)
	//req.Header.Add()
	//
	//fasthttp.DoTimeout(req, resp, time.Second)
	//
	//bodyBytes := resp.Body()
	//println(string(bodyBytes))
	//// User-Agent: fasthttp
	//// Body:
	//
	//fasthttp.Post()
	//
	//client := fasthttp.HostClient{}
	//client.DoTimeout()
	//
	//fasthttp.Client{
	//	Dial:
	//}
}
