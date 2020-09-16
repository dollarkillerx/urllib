package urllib

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func Get(url string) *Urllib {
	return getBase(url, "GET")
}

func Post(url string) *Urllib {
	return getBase(url, "POST")
}

func Put(url string) *Urllib {
	return getBase(url, "PUT")
}

func Delete(url string) *Urllib {
	return getBase(url, "DELETE")
}

func Head(url string) *Urllib {
	return getBase(url, "HEAD")
}

// 公开Urllib 方便未来定制包装
type Urllib struct {
	url string

	req  *http.Request
	resp *http.Response

	params url.Values
	querys url.Values

	data   []byte
	files  map[string]string
	config Config

	err error
}

type Config struct {
	UserAgent        string
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	TLSClientConfig  *tls.Config
	Proxy            func(*http.Request) (*url.URL, error)
	Transport        http.RoundTripper
	CheckRedirect    func(req *http.Request, via []*http.Request) error
	Gzip             bool
}

var defaultConfig = Config{
	UserAgent:        "Urllib Gold",
	ConnectTimeout:   30 * time.Second,
	ReadWriteTimeout: 30 * time.Second,
	Gzip:             false,
}

func getBase(tagUrl, method string) *Urllib {
	var resp http.Response
	u, err := url.Parse(tagUrl)

	base := &Urllib{
		params: map[string][]string{},
		querys: map[string][]string{},
		files:  map[string]string{},
		url:    tagUrl,
		err:    err,
		config: defaultConfig,
		resp:   &resp,
	}

	base.req = &http.Request{
		URL:        u,
		Method:     method,
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	base.SetUserAgent(base.config.UserAgent)
	return base
}

// 允许所有TLS
func (u *Urllib) AlloverTLS() *Urllib {
	u.config.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return u
}

func (u *Urllib) SetGzip() *Urllib {
	u.config.Gzip = true
	return u
}

func (u *Urllib) SetConfig(config Config) *Urllib {
	u.config = config
	return u
}

func (u *Urllib) SetHeader(k, v string) *Urllib {
	u.req.Header.Set(k, v)
	return u
}

func (u *Urllib) SetHeaderMap(headMap map[string]string) *Urllib {
	for k, v := range headMap {
		u.req.Header.Set(k, v)
	}
	return u
}

func (u *Urllib) SetHost(host string) *Urllib {
	u.req.Host = host
	return u
}

func (u *Urllib) SetUserAgent(ua string) *Urllib {
	u.req.Header.Set("User-Agent", ua)
	return u
}

func (u *Urllib) RandUserAgent() *Urllib {
	u.req.Header.Set("User-Agent", ReptileGetUserAgent())
	return u
}

func (u *Urllib) Params(key, val string) *Urllib {
	u.params.Add(key, val)
	return u
}

func (u *Urllib) Queries(key, val string) *Urllib {
	u.querys.Add(key, val)
	return u
}

func (u *Urllib) ParamsMap(data map[string]string) *Urllib {
	for k, v := range data {
		u.params.Add(k, v)
	}
	return u
}

func (u *Urllib) QueriesMap(data map[string]string) *Urllib {
	for k, v := range data {
		u.querys.Add(k, v)
	}
	return u
}

func (u *Urllib) SetTimeout(timeout time.Duration) *Urllib {
	if timeout < 10 {
		timeout = timeout * time.Second
	}
	u.config.ConnectTimeout = timeout
	u.config.ReadWriteTimeout = timeout
	return u
}

func (u *Urllib) SetConnectTimeout(timeout time.Duration) *Urllib {
	if timeout < 10 {
		timeout = timeout * time.Second
	}
	u.config.ConnectTimeout = timeout
	return u
}

func (u *Urllib) SetReadWriteTimeout(timeout time.Duration) *Urllib {
	if timeout < 10 {
		timeout = timeout * time.Second
	}
	u.config.ReadWriteTimeout = timeout
	return u
}

func (u *Urllib) SetTLSClientConfig(config *tls.Config) *Urllib {
	u.config.TLSClientConfig = config
	return u
}

func (u *Urllib) SetHTTPVersion(version string) *Urllib {
	major, minor, ok := http.ParseHTTPVersion(version)
	if ok {
		u.req.Proto = version
		u.req.ProtoMajor = major
		u.req.ProtoMinor = minor
	}

	return u
}

func (u *Urllib) SetCookie(cookie *http.Cookie) *Urllib {
	u.req.Header.Add("Cookie", cookie.String())
	return u
}

func (u *Urllib) SetTransport(transport http.RoundTripper) *Urllib {
	u.config.Transport = transport
	return u
}

func (u *Urllib) SetCheckRedirect(redirect func(req *http.Request, via []*http.Request) error) *Urllib {
	u.config.CheckRedirect = redirect
	return u
}

func (u *Urllib) SetAuth(user, password string) *Urllib {
	u.req.SetBasicAuth(user, password)
	return u
}

func (u *Urllib) SetProxy(proxy func(*http.Request) (*url.URL, error)) *Urllib {
	u.config.Proxy = proxy
	return u
}

func (u *Urllib) PostFile(fromKey, filePath string) *Urllib {
	u.files[fromKey] = filePath
	return u
}

func (u *Urllib) SetBody(body []byte) *Urllib {
	u.req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	u.req.ContentLength = int64(len(body))
	return u
}

func (u *Urllib) SetJson(body []byte) *Urllib {
	if u.req.Body == nil && body != nil {
		u.req.Body = ioutil.NopCloser(bytes.NewReader(body))
		u.req.ContentLength = int64(len(body))
		u.req.Header.Set("Content-Type", "application/json")
	}
	return u
}

func (u *Urllib) SetJsonObject(obj interface{}) *Urllib {
	if u.req.Body == nil && obj != nil {
		body, err := json.Marshal(obj)
		if err != nil {
			u.err = err
			return u
		}
		u.req.Body = ioutil.NopCloser(bytes.NewReader(body))
		u.req.ContentLength = int64(len(body))
		u.req.Header.Set("Content-Type", "application/json")
	}
	return u
}

// if form
func (u *Urllib) setFormParams(Forms url.Values) {
	data := Forms.Encode()
	u.req.Body = ioutil.NopCloser(strings.NewReader(data))
	u.req.ContentLength = int64(len(data))
}

func (u *Urllib) setGzip() *Urllib {
	if u.req.Body == nil {
		return u
	}

	if u.req.Body == nil {
		return u
	}
	all, err := ioutil.ReadAll(u.req.Body)
	if err != nil {
		u.err = err
		return u
	}
	if len(all) == 0 {
		return u
	}
	data, err := GZipData(all)
	if err != nil {
		u.err = err
		return u
	}

	u.SetHeader("Content-Encoding", "gzip")
	u.req.Body = ioutil.NopCloser(bytes.NewReader(data))
	u.req.ContentLength = int64(len(data))
	return u
}

func (u *Urllib) NoRedirect() *Urllib {
	u.config.CheckRedirect = func(req *http.Request, via []*http.Request) error { // 定制禁用跳转
		return http.ErrUseLastResponse
	}
	return u
}

func (u *Urllib) basicRules() {
	// Url Params
	distUrl, err := buildURLParams(u.url, u.querys)
	if err != nil {
		u.err = err
		return
	}
	parse, err := url.Parse(distUrl)
	if err != nil {
		u.err = err
		return
	}
	u.req.URL = parse

	if len(u.params) > 0 {
		u.setFormParams(u.params)
		u.SetHeader("Content-Type", "application/x-www-form-urlencoded")
		return
	}

	if len(u.files) > 0 {
		pr, pw := io.Pipe()
		bodyWriter := multipart.NewWriter(pw)
		go func() {
			for formname, filename := range u.files {
				fileWriter, err := bodyWriter.CreateFormFile(formname, filename)
				if err != nil {
					log.Println("Urllib:", err)
				}
				fh, err := os.Open(filename)
				if err != nil {
					log.Println("Urllib:", err)
				}
				//iocopy
				_, err = io.Copy(fileWriter, fh)
				fh.Close()
				if err != nil {
					log.Println("Urllib:", err)
				}
			}
			for k, v := range u.params {
				for _, vv := range v {
					bodyWriter.WriteField(k, vv)
				}
			}
			bodyWriter.Close()
			pw.Close()
		}()
		u.SetHeader("Content-Type", bodyWriter.FormDataContentType())
		u.req.Body = ioutil.NopCloser(pr)
		u.SetHeader("Transfer-Encoding", "chunked")
		return
	}
}

func (u *Urllib) body() (*http.Response, error) {
	u.basicRules()

	if u.config.Gzip {
		u.setGzip()
	}

	if u.config.Transport == nil {
		u.config.Transport = &http.Transport{
			TLSClientConfig:     u.config.TLSClientConfig,
			Proxy:               u.config.Proxy,
			Dial:                setTimeoutDialer(u.config.ConnectTimeout, u.config.ReadWriteTimeout),
			MaxIdleConnsPerHost: 100,
		}
	} else {
		if t, ok := u.config.Transport.(*http.Transport); ok {
			if t.TLSClientConfig == nil {
				t.TLSClientConfig = u.config.TLSClientConfig
			}
			if t.Proxy == nil {
				t.Proxy = u.config.Proxy
			}
			if t.Dial == nil {
				t.Dial = setTimeoutDialer(u.config.ConnectTimeout, u.config.ReadWriteTimeout)
			}
		}
	}

	initCookie()
	jar := globalCookie

	client := http.Client{
		Transport: u.config.Transport,
		Jar:       jar,
	}

	if u.config.CheckRedirect != nil {
		client.CheckRedirect = u.config.CheckRedirect
	}

	return client.Do(u.req)
}

func (u *Urllib) byte() (int, []byte, error) {
	body, err := u.body()
	if err != nil {
		return 0, nil, err
	}
	defer body.Body.Close()
	all, err := ioutil.ReadAll(body.Body)
	if err != nil {
		return 0, all, err
	}
	// GZIP
	if body.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(body.Body)
		if err != nil {
			return 0, nil, err
		}
		all, err = ioutil.ReadAll(reader)
		if err != nil {
			return 0, nil, err
		}
	}

	// 旋转木马
	contentType := body.Header.Get("Content-Type")
	all, err = AutomaticTranscoding(contentType, all)
	return body.StatusCode, all, err
}

func (u *Urllib) Body() (*http.Response, error) {
	return u.body()
}
func (u *Urllib) Byte() (int, []byte, error) {
	return u.byte()
}

// 拥有重尝 版本
func (u *Urllib) BodyRetry(retry int) (body *http.Response, err error) {
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

func (u *Urllib) ByteRetry(retry int) (statusCode int, body []byte, err error) {
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

var globalCookie http.CookieJar
var cookieMutex sync.Mutex

func initCookie() {
	cookieMutex.Lock()
	defer cookieMutex.Unlock()
	if globalCookie == nil {
		jar, err := cookiejar.New(nil)
		if err == nil {
			globalCookie = jar
			return
		}
		log.Println(err)
	}
}
