package test

import (
	"fmt"
	"github.com/dollarkillerx/urllib"
	"github.com/dollarkillerx/urllib/lib"
	"io/ioutil"
	"log"
	"testing"
)

func TestSendJson(t *testing.T) {
	//i, bytes, err := urllib.Post("http://0.0.0.0:8985/test_body").
	i, bytes, err := urllib.Post("https://www.baidu.com").
		SetHeader("AUTH", "b095f33d75a248da9255e822ffe859aa").
		SetJson([]byte("Hello World")).ByteRetry(3)
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
}
