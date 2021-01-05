package test

import (
	"crypto/tls"
	"fmt"
	"github.com/dollarkillerx/fasthttp"
	"github.com/dollarkillerx/urllib"
	"github.com/dollarkillerx/urllib/lib"
	"log"
	"testing"
)

//func TestSendJson(t *testing.T) {
//	//i, bytes, err := urllib.Post("http://0.0.0.0:8985/test_body").
//	i, bytes, err := urllib.Post("https://www.baidu.com").
//		SetHeader("AUTH", "b095f33d75a248da9255e822ffe859aa").
//		SetJson([]byte("Hello World")).ByteRetry(3)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Println(i)
//	log.Println(string(bytes))
//}

func TestRedirect(t *testing.T) {
	get, body, err := fasthttp.Get(nil, "http://www.baidu.com")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(get)
	log.Println(string(body))

	//fasthttp.Post()
}

func TestRedirect2(t *testing.T) {
	req := &fasthttp.Request{}
	req.SetRequestURI("http://127.0.0.1:8986/t")
	req.Header.SetMethod("GET")
	resp := &fasthttp.Response{}

	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		log.Fatalln(err)
	}
	log.Println(resp.StatusCode())
}

func TestHttps(t *testing.T) {
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	//client := fasthttp.HostClient{
	//	IsTLS:     true,
	//	TLSConfig: tlsConfig,
	//}
	client := fasthttp.Client{
		TLSConfig: tlsConfig,
	}
	//req.URI().SetScheme("https")
	req.SetRequestURI("https://0.0.0.0:8091/")
	//req.Header.SetHost("httpbin.org")
	err := client.Do(req, resp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(resp.Body()))
}

func TestArgs(t *testing.T) {
	a := &fasthttp.Args{}
	//a.Add("asd","sad")
	log.Println(a.String())
}

func TestGet(t *testing.T) {
	i, bytes, err := urllib.Get("http://127.0.0.1:8986/test").Byte()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(i)
	log.Println(string(bytes))
}


func TestMokeIP(t *testing.T) {
	for i:=0;i<10;i++ {
		fmt.Println(lib.RandomIp())
	}
}