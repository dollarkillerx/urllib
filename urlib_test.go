package urllib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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


func TestAutomaticTranscoding(t *testing.T) {
	i, bytes, err := Get("https://www.discuz.net/forum.php").Byte()  // gbk
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(bytes))
	fmt.Println(i)

	i, bytes, err = Get("http://www.phome.net/").Byte()             // 无标注
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(bytes))
	fmt.Println(i)
}

func TestPost2(t *testing.T) {
	body, err := Get("http://127.0.0.1:8082/").Body()
	if err != nil {
		log.Fatalln(err)
	}
	defer body.Body.Close()
	all, err := ioutil.ReadAll(body.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(all))

	fmt.Println(body.StatusCode)

	resp, err := http.Get("http://127.0.0.1:8082/")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(readAll))
	log.Println(resp.StatusCode)
}