package urllib

import (
	"fmt"
	"log"
	"net/url"
	"testing"
)

func TestLib(t *testing.T) {
	ur := "http://baidu.com/adssd/sadsad"
	u := url.Values{}
	u.Add("user", "root")
	u.Add("password", "pwd")

	params, err := buildURLParams(ur, u)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(params)
}

func TestGet(t *testing.T) {
	body, err := Get("http://www.baidu.com").Body()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(body.StatusCode)

	//bod, err := Get("http://www.baidu.com").Byte()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(string(bod))
}

func TestPost(t *testing.T) {
	post := Post("http://168.1xxxxx/cg")
	post.Params("username", "root")
	post.Params("password", "we")
	body, err := post.Body()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(body.StatusCode)
}

func TestTimeOut(t *testing.T) {
	//bytes, err := Get("http://www.google.com").
	//	HttpProxy("http://proxy.com").
	//	SetTimeout(time.Second * 3).
	//	Byte()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(bytes)
}

func TestGetQuery(t *testing.T) {
	retry, body, err := Get("http://192.168.89.56:8080/assets").Params("url", "test.com/2.html").Params("url", "test.com/1.html").ByteRetry(3)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(retry)
	fmt.Println(string(body))
}
