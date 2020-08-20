package main

import (
	"fmt"
	"github.com/dollarkillerx/urllib"
	"log"
)

func main() {
	// get
	httpCode, bytes, err := urllib.Get("http://www.baidu.com").Byte()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("HttpCode: %d \n", httpCode)
	fmt.Println(string(bytes))


	httpCode, bytes, err = urllib.Get("http://www.baidu.com").
		Query("q","122").
		Query("h","1213").Byte()   // 生成URL： http://www.baidu.com?q=122&h=1213
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("HttpCode: %d \n", httpCode)
	fmt.Println(string(bytes))


	httpCode, bytes, err = urllib.Post("http://www.baidu.com").
		Params("q","122").
		Params("h","1213").Byte()   // URL： http://www.baidu.com  表单 q=122 h=1213
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("HttpCode: %d \n", httpCode)
	fmt.Println(string(bytes))

	httpCode, bytes, err = urllib.Post("http://www.baidu.com").
		SetJson([]byte(`{"MSG":"121321"}`)).Byte()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("HttpCode: %d \n", httpCode)
	fmt.Println(string(bytes))


	httpCode, bytes, err = urllib.Post("http://www.baidu.com").HttpProxy("http://xxxx.c").Byte()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("HttpCode: %d \n", httpCode)
	fmt.Println(string(bytes))

	httpCode, bytes, err = urllib.Post("http://www.baidu.com").SetTimeout(3).Byte()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("HttpCode: %d \n", httpCode)
	fmt.Println(string(bytes))


	httpCode, bytes, err = urllib.Delete("http://www.baidu.com").ByteRetry(3)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("HttpCode: %d \n", httpCode)
	fmt.Println(string(bytes))

	body, err := urllib.Post("http://www.baidu.com").
		SetHeader("xxx", "xxx").
		SetHeader("xxx", "xxx").Body()
	defer body.Body.Close()

	urllib.Post("http://www.baidu.com").SetAuth("user","passwd").Body()

	urllib.Post("http://www.baidu.com").SetUserAgent("chrome").Body()
	urllib.Post("http://www.baidu.com").RandUserAgent().Body()

	//urllib.Post("http://www.baidu.com").SetCookie().Body()
}
