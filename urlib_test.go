package urllib

import (
	"io/ioutil"
	"log"
	"testing"
)

//func TestLib(t *testing.T) {
//	ur := "http://baidu.com/adssd/sadsad"
//	u := url.Values{}
//	u.Add("user", "root")
//	u.Add("password", "pwd")
//
//	params, err := buildURLParams(ur, u)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Println(params)
//}
//
//func TestGet(t *testing.T) {
//	body, err := Get("http://www.baidu.com").Body()
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Println(body.StatusCode)
//
//	//bod, err := Get("http://www.baidu.com").Byte()
//	//if err != nil {
//	//	log.Fatalln(err)
//	//}
//	//fmt.Println(string(bod))
//}
//
//func TestPost(t *testing.T) {
//	post := Post("http://168.1xxxxx/cg")
//	post.Params("username", "root")
//	post.Params("password", "we")
//	body, err := post.Body()
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Println(body.StatusCode)
//}
//
//func TestTimeOut(t *testing.T) {
//	//bytes, err := Get("http://www.google.com").
//	//	HttpProxy("http://proxy.com").
//	//	SetTimeout(time.Second * 3).
//	//	Byte()
//	//if err != nil {
//	//	log.Fatalln(err)
//	//}
//	//fmt.Println(bytes)
//}
//
//func TestGetQuery(t *testing.T) {
//	retry, body, err := Get("http://192.168.89.56:8080/assets").Params("url", "test.com/2.html").Params("url", "test.com/1.html").ByteRetry(3)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Println(retry)
//	fmt.Println(string(body))
//}
//
//
//func TestAutomaticTranscoding(t *testing.T) {
//	i, bytes, err := Get("https://www.discuz.net/forum.php").Byte()  // gbk
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Println(string(bytes))
//	fmt.Println(i)
//
//	i, bytes, err = Get("http://www.phome.net/").Byte()             // 无标注
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Println(string(bytes))
//	fmt.Println(i)
//}
//
//func TestPost2(t *testing.T) {
//	body, err := Get("http://127.0.0.1:8082/").Body()
//	if err != nil {
//		log.Fatalln(err)
//	}
//	defer body.Body.Close()
//	all, err := ioutil.ReadAll(body.Body)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Println(string(all))
//
//	fmt.Println(body.StatusCode)
//
//	resp, err := http.Get("http://127.0.0.1:8082/")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	defer resp.Body.Close()
//	readAll, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Println(string(readAll))
//	log.Println(resp.StatusCode)
//}

func TestAbc(t *testing.T) {
	//urls := "http://192.168.88.11:9001/assets/2083e95e-7928-4292-845c-a4c794a23282.html?" +
	//	"X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=minio%2F20200831%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20200831T084536Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=4f6e3f7668fdff84eb7bfdb1bdbfd25dc6057d330248ed2e37aeaca89efea56f"
	//parse, err := url.Parse(urls)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//urls = parse.String()
	//log.Println(urls)
	//
	//abs := ""
	//sp := make(map[string]string)
	//if index := strings.Index(urls, "?"); index > 0 {
	//	abs = urls[index+1:]
	//}
	//if abs != "" {
	//	split := strings.Split(abs, "&")
	//	for _,k := range split {
	//		pcr := strings.Split(k, "=")
	//		if len(pcr) == 2 {
	//			sp[pcr[0]] = pcr[1]
	//		}else {
	//			sp[pcr[0]] = ""
	//		}
	//	}
	//}
	//
	//retry, body, err := Get("http://192.168.88.11:9001/assets/2083e95e-7928-4292-845c-a4c794a23282.html").QueriesMap(sp).ByteRetry(3)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Println(string(body))
	//log.Println(retry)
	//
	//byteRetry, bytes, err := Get("http://192.168.88.11:9001/assets/2083e95e-7928-4292-845c-a4c794a23282.html").
	//	QueriesMap(map[string]string{
	//		"X-Amz-Algorithm": "AWS4-HMAC-SHA256",
	//		"X-Amz-Credential": "minio%2F20200831%2Fus-east-1%2Fs3%2Faws4_request",
	//		"X-Amz-Date": "20200831T084536Z",
	//		"X-Amz-Expires": "1800",
	//		"X-Amz-SignedHeaders": "host",
	//		"X-Amz-Signature": "4f6e3f7668fdff84eb7bfdb1bdbfd25dc6057d330248ed2e37aeaca89efea56f",
	//}).ByteRetry(3)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//log.Println(string(bytes))
	//log.Println(byteRetry)

	//resp, err := http.Get("http://192.168.88.11:9001/assets/2083e95e-7928-4292-845c-a4c794a23282.html?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=minio%2F20200831%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20200831T095914Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=ccd5387da9dbcb406b2b9ef250d2a619fbdd323f31fb7e47dd44e036fa108805")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//all, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Println(resp.Header)
	//log.Println(resp.StatusCode)
	//log.Println(string(all))

	url := "http://192.168.88.11:9001/assets/2083e95e-7928-4292-845c-a4c794a23282.html?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=minio%2F20200831%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20200831T095914Z&X-Amz-Expires=1800&X-Amz-SignedHeaders=host&X-Amz-Signature=ccd5387da9dbcb406b2b9ef250d2a619fbdd323f31fb7e47dd44e036fa108805"
	retry, body, err := Get(url).ByteRetry(3)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
	log.Println(retry)
}
func TestAbb(t *testing.T) {
	Get("http://www.baidu.com?q=1212").ByteRetry(3)
}

func TestPost(t *testing.T) {
	retry, body, err := Post("http://0.0.0.0:8083/test").SetJsonObject(map[string]interface{}{
		"name": "dolalark",
		"pc": map[string]string{
			"ads": "sada",
		},
	}).ByteRetry(3)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(retry)
	log.Println(string(body))
}

func TestPost2(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	body, err := Post("http://0.0.0.0:8083/test").SetJsonObject(map[string]interface{}{
		"name": "dolalark",
		"pc": map[string]string{
			"ads": "sada",
		},
	}).Body()
	if err != nil {
		log.Fatalln(err)
	}
	defer body.Body.Close()
	log.Println(body.StatusCode)

	all, err := ioutil.ReadAll(body.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(all))
}


func TestPostFrom(t *testing.T) {
	retry, body, err := Post("http://0.0.0.0:8986/tc").Params("Name", "sadas").Params("Age", "16").ByteRetry(3)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(retry)
	log.Println(string(body))
}

func TestPostJson(t *testing.T) {
	type p struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	retry, body, err := Post("http://0.0.0.0:8986/test").SetJsonObject(&p{
		Name: "dollarkiller",
		Age:  20,
	}).ByteRetry(3)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(retry)
	log.Println(string(body))
}

func TestGetQuery(t *testing.T) {
	retry, body, err := Get("http://0.0.0.0:8986/c").Queries("q", "hello").Queries("g","asdsad").ByteRetry(3)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(retry)
	log.Println(string(body))

	retry, body, err = Get("http://0.0.0.0:8986/c?name=sdsd&q=sad").Queries("g","asdsad").ByteRetry(3)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(retry)
	log.Println(string(body))
}

func TestGet(t *testing.T) {
	i, bytes, err := Get("http://127.0.0.1:8986/test").Byte()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(i)
	log.Println(string(bytes))

	i, bytes, err = Get("http://127.0.0.1:8986/t").NoRedirect().Byte()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(i)
	log.Println(string(bytes))
}


func TestPojs(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	js := `{"partner_id":"web","param":{"pageCallback":"eyJmaXJzdElkIjo4OTk5ODU3NTM4NjI5MTIsImxhc3RJZCI6ODk5OTE3MDg5OTQ0ODM3LCJmaXJzdENyZWF0ZVRpbWUiOjE2MDEyMDI3NzExMDksImxhc3RDcmVhdGVUaW1lIjoxNjAxMTk4NTgwMTk2fQ","pageEvent":1,"pageSize":20,"platformId":2,"siteId":1}}`

	//retry, body, err := urllib.Post("https://gateway.36kr.com/api/mis/nav/newsflash/flow").
	retry, body, err := Post("http://127.0.0.1:8081/test").
		SetJson([]byte(js)).SetHeader("Content-Type","application/json;charset=UTF-8").ByteRetry(3)
	if err != nil {
		log.Println("a")
		log.Fatalln(err)
	}
	log.Println(retry)
	log.Println(string(body))
}
func TestGp(t *testing.T)  {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	i, bytes, err := Get("https://36kr.com/newsflashes").Byte()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(i)
	log.Println(string(bytes))
}