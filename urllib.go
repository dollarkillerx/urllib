package urllib

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type urlType string

const (
	UA = "DollarKiller-UrlLib2.0"

	get    urlType = "GET"
	post   urlType = "POST"
	put    urlType = "PUT"
	delete urlType = "DELETE"
)

func Get(url string) *urllib {
	base := getBase()
	base.url = url
	base.typ = get
	return base
}

func Post(url string) *urllib {
	base := getBase()
	base.url = url
	base.typ = post
	base.req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return base
}

func Put(url string) *urllib {
	base := getBase()
	base.url = url
	base.typ = put
	return base
}

func Delete(url string) *urllib {
	base := getBase()
	base.url = url
	base.typ = delete
	return base
}

type urllib struct {
	url     string
	typ     urlType
	timeout time.Duration

	noTls     bool
	httpProxy string
	req       *http.Request
	resp      *http.Response

	params  url.Values
	querys  url.Values
	header  map[string]string
	cookies []*http.Cookie

	client *http.Client

	bodyByte []byte
}

func getBase() *urllib {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	base := &urllib{
		timeout: time.Second * 3,
		params:  map[string][]string{},
		querys:  map[string][]string{},
		header:  map[string]string{},
		cookies: []*http.Cookie{},
		client:  &http.Client{},
	}

	base.client.Timeout = time.Second * 3
	base.client.Jar = jar

	base.req = &http.Request{
		Method:     string(get),
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	base.req.Header.Set("User-Agent", UA)
	return base
}

func (u *urllib) SetUserAgent(ua string) *urllib {
	u.req.Header.Set("User-Agent", ua)
	return u
}

func (u *urllib) RandUserAgent() *urllib {
	u.req.Header.Set("User-Agent", ReptileGetUserAgent())
	return u
}

func (u *urllib) SetJson(body []byte) *urllib {
	u.SetHeader("Content-Type", "application/json;charset=UTF-8")
	u.bodyByte = body
	return u
}

func (u *urllib) Params(key, val string) *urllib {
	u.params.Add(key, val)
	return u
}

func (u *urllib) Queries(key, val string) *urllib {
	u.querys.Add(key, val)
	return u
}

func (u *urllib) ParamsMap(data map[string]string) *urllib {
	for k, v := range data {
		u.params.Add(k, v)
	}
	return u
}

func (u *urllib) QueriesMap(data map[string]string) *urllib {
	for k, v := range data {
		u.querys.Add(k, v)
	}
	return u
}

func (u *urllib) SetCookie(cookies []*http.Cookie) *urllib {
	u.cookies = cookies
	return u
}

func (u *urllib) SetHeader(k, v string) *urllib {
	u.header[k] = v
	return u
}

func (u *urllib) SetHeaderMap(headMap map[string]string) *urllib {
	for k, v := range headMap {
		u.header[k] = v
	}
	return u
}

// 不检测TLS
func (u *urllib) NoTLS() *urllib {
	u.noTls = true
	return u
}

func (u *urllib) HttpProxy(url string) *urllib {
	u.httpProxy = url
	return u
}

func (u *urllib) ClearCookies() {
	u.cookies = u.cookies[0:0]
}

func (u *urllib) clientSetCookies() {
	if len(u.cookies) > 0 {
		u.client.Jar.SetCookies(u.req.URL, u.cookies)
		u.ClearCookies()
	}
}

func (u *urllib) SetAuth(user, password string) *urllib {
	u.req.SetBasicAuth(user, password)
	return u
}

func (u *urllib) SetTimeout(n time.Duration) *urllib {
	if n < 10 {
		n = n * time.Second
	}
	u.client.Timeout = n
	return u
}

func (u *urllib) SetGzip() *urllib {
	u.SetHeader("Content-Encoding", "gzip")
	return u
}

func (u *urllib) SetBody(body []byte) *urllib {
	u.bodyByte = body
	return u
}

func (u *urllib) body() (*http.Response, error) {
	var baseUrl *url.URL
	// querys ?xx=xx&xx=xx
	u.setBodyBytes(u.params)
	disturl, err := buildURLParams(u.url, u.querys)
	if err != nil {
		return nil, err
	}
	parse, err := url.Parse(disturl)
	if err != nil {
		return nil, err
	}
	baseUrl = parse
	// querys end

	switch u.typ {
	case get:
		u.req.Method = string(get)
	case post:
		u.req.Method = string(post)
	case put:
		u.req.Method = string(put)
	case delete:
		u.req.Method = string(delete)
	}

	// gip
	if u.header["Content-Encoding"] == "gzip" {
		data, err := GZipData(u.bodyByte)
		if err != nil {
			return nil, err
		}
		u.bodyByte = data
	}

	// set json
	if u.header["Content-Type"] == "application/json;charset=UTF-8" {
		request, err := http.NewRequest("POST", u.url, bytes.NewBuffer(u.bodyByte))
		if err != nil {
			return nil, err
		}
		u.req = request
	}

	u.req.URL = baseUrl
	u.clientSetCookies()

	for k, v := range u.header {
		u.req.Header.Add(k, v)
	}

	if u.httpProxy != "" {
		urli := url.URL{}
		urlproxy, err := urli.Parse(u.httpProxy)
		if err != nil {
			fmt.Println("Set proxy failed")
			return nil, err
		}
		u.client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(urlproxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !u.noTls},
		}
	} else {
		u.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !u.noTls},
		}
	}

	return u.client.Do(u.req)
}

func (u *urllib) byte() (int, []byte, error) {
	body, err := u.body()
	if err != nil {
		return 0, nil, err
	}
	defer body.Body.Close()
	all, err := ioutil.ReadAll(body.Body)
	if err != nil {
		return 0, all, err
	}
	contentType := body.Header.Get("Content-Type")
	all, err = AutomaticTranscoding(contentType, all)
	return body.StatusCode, all, err
}

func (u *urllib) setBodyBytes(Forms url.Values) {
	data := Forms.Encode()
	u.req.Body = ioutil.NopCloser(strings.NewReader(data))
	u.req.ContentLength = int64(len(data))
}

func (u *urllib) NoRedirect() *urllib {
	u.client.CheckRedirect = func(req *http.Request, via []*http.Request) error { // 定制禁用跳转
		return http.ErrUseLastResponse
	}
	return u
}

func (u *urllib) Body() (*http.Response, error) {
	return u.body()
}
func (u *urllib) Byte() (int, []byte, error) {
	return u.byte()
}

// 拥有重尝 版本
func (u *urllib) BodyRetry(retry int) (body *http.Response, err error) {
	for i := 0; i < retry; i++ {
		body, err = u.body()
		if err != nil {
			time.Sleep(time.Second * time.Duration(random(1, 5)))
			continue
		}
		return body, err
	}
	return body, err
}

func (u *urllib) ByteRetry(retry int) (statusCode int, body []byte, err error) {
	if retry == 0 {
		retry = 3
	}
	for i := 0; i < retry; i++ {
		statusCode, body, err = u.byte()
		if err != nil {
			time.Sleep(time.Second * time.Duration(random(1, 5)))
			continue
		}
		return statusCode, body, err
	}
	return statusCode, body, err
}
