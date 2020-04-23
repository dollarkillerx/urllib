package urllib

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type urlType string

const (
	UA = "UrlLib By DollarKillerx"

	GET    urlType = "GET"
	POST   urlType = "POST"
	PUT    urlType = "PUT"
	DELETE urlType = "DELETE"
)

type urllib struct {
	url     string
	typ     urlType
	timeout time.Duration

	noTls     bool
	httpProxy string
	req       *http.Request
	resp      *http.Response

	params  url.Values
	header  map[string]string
	cookies []*http.Cookie

	client *http.Client
}

func getBase() *urllib {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	base := &urllib{
		timeout: time.Second * 3,
		params:  map[string][]string{},
		header:  map[string]string{},
		cookies: []*http.Cookie{},
		client:  &http.Client{},
	}

	base.client.Jar = jar

	base.req = &http.Request{
		Method:     string(GET),
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	base.req.Header.Set("User-Agent", UA)
	return base
}

func Get(url string) *urllib {
	base := getBase()
	base.url = url
	base.typ = GET
	return base
}

func Post(url string) *urllib {
	base := getBase()
	base.url = url
	base.typ = POST
	base.req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return base
}

func Put(url string) *urllib {
	base := getBase()
	base.url = url
	base.typ = PUT
	return base
}

func Delete(url string) *urllib {
	base := getBase()
	base.url = url
	base.typ = DELETE
	return base
}

func (u *urllib) SetTimeOut(timeout time.Duration) *urllib {
	u.timeout = timeout
	return u
}

func (u *urllib) SetUserAgent(ua string) *urllib {
	u.req.Header.Set("User-Agent", ua)
	return u
}

func (u *urllib) RandUserAgent() *urllib {
	u.req.Header.Set("User-Agent", ReptileGetUserAgent())
	return u
}

func (u *urllib) Params(key, val string) *urllib {
	u.params.Set(key, val)
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
		// 1. Cookies have content, Copy Cookies to Client.jar
		// 2. Clear  Cookies
		u.client.Jar.SetCookies(u.req.URL, u.cookies)
		u.ClearCookies()
	}
}

func (u *urllib) SetAuth(user, password string) *urllib {
	u.req.SetBasicAuth(user, password)
	return u
}

func (u *urllib) SetTimeout(n time.Duration) *urllib {
	u.client.Timeout = n
	return u
}

func (u *urllib) Body() (*http.Response, error) {
	delete(u.req.Header, "Cookie")
	var baseUrl *url.URL
	switch u.typ {
	case GET:
		disturl, _ := buildURLParams(u.url, u.params)
		parse, err := url.Parse(disturl)
		if err != nil {
			return nil, err
		}
		baseUrl = parse
	case POST:
		u.setBodyBytes(u.params)
		parse, err := url.Parse(u.url)
		if err != nil {
			return nil, err
		}
		baseUrl = parse
	}

	u.req.URL = baseUrl
	u.clientSetCookies()

	for k, v := range u.header {
		u.req.Header.Set(k, v)
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

func (u *urllib) Byte() ([]byte, error) {
	body, err := u.Body()
	if err != nil {
		return nil, err
	}
	defer body.Body.Close()
	return ioutil.ReadAll(body.Body)
}
