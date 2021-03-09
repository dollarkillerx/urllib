package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/dollarkillerx/urllib"
	"github.com/dollarkillerx/urllib/lib"
)

func TestSendJson(t *testing.T) {
	//i, bytes, err := urllib.Post("http://0.0.0.0:8985/test_body").
	i, bytes, err := urllib.Post("https://www.baidu.com").
		SetHeader("AUTH", "b095f33d75a248da9255e822ffe859aa").
		SetJson([]byte("Hello World")).ByteRetry(3, 200)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(i)
	log.Println(string(bytes))
}

func TestIp(t *testing.T) {
	fmt.Println(lib.RandomIp())
}

func TestSendIp(t *testing.T) {
	_, bytes, err := urllib.Get("https://www.ez2o.com/App/Net/IP").RandDisguisedIP().Byte()
	if err != nil {
		log.Fatalln(err)
	}
	ioutil.WriteFile("ip.html", bytes, 00666)

	_, bytes, err = urllib.Get("https://www.ez2o.com/App/Net/IP").DisguisedIP("1.1.1.1").Byte()
	if err != nil {
		log.Fatalln(err)
	}
	ioutil.WriteFile("ip2.html", bytes, 00666)
}

func TestP2(t *testing.T) {
	r, body, err := urllib.Get("https://news.windin.com/ns/imagebase/6316/6683.9fa2d2eb36f1c3940856f104ec77fc98").
		SetTimeout(time.Second).RandUserAgent().KeepAlives().Byte()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(r)
	ioutil.WriteFile("a.png", body, 00666)
}
