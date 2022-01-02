package urllib

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dollarkillerx/urllib/lib"
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

	disableKeepAlives bool

	err error
}

type Config struct {
	UserAgent       string
	ConnectTimeout  time.Duration
	TLSClientConfig *tls.Config
	Proxy           func(*http.Request) (*url.URL, error)
	Transport       http.RoundTripper
	CheckRedirect   func(req *http.Request, via []*http.Request) error
	Gzip            bool
	Debug           bool
}

var defaultConfig = Config{
	UserAgent:      "Urllib Gold",
	ConnectTimeout: 30 * time.Second,
	Gzip:           false,
}

func getBase(tagUrl, method string) *Urllib {
	tagUrl = lib.AddPlt(tagUrl)
	var resp http.Response
	u, err := url.Parse(tagUrl)

	base := &Urllib{
		params:            map[string][]string{},
		querys:            map[string][]string{},
		files:             map[string]string{},
		url:               tagUrl,
		err:               err,
		config:            defaultConfig,
		resp:              &resp,
		disableKeepAlives: true,
	}

	base.req = &http.Request{
		URL:        u,
		Method:     method,
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Close:      true,
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
	u.req.Header.Set("User-Agent", lib.ReptileGetUserAgent())
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

func (u *Urllib) setBody() {
	if len(u.data) != 0 {
		u.req.Body = ioutil.NopCloser(bytes.NewBuffer(u.data))
		u.req.ContentLength = int64(len(u.data))
	}
}

func (u *Urllib) SetBody(body []byte) *Urllib {
	u.data = body
	return u
}

func (u *Urllib) SetJson(body []byte) *Urllib {
	if body != nil {
		u.data = body
		u.req.Header.Set("Content-Type", "application/json")
	}
	return u
}

// 防止请求服务器Socks关闭引起 EOF错误
func (u *Urllib) PreventEOF() *Urllib {
	u.req.Close = true
	return u
}

func (u *Urllib) SetJsonObject(obj interface{}) *Urllib {
	if obj != nil {
		body, err := json.Marshal(obj)
		if err != nil {
			u.err = err
			return u
		}
		u.data = body
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
	data, err := lib.GZipData(all)
	if err != nil {
		u.err = err
		if u.config.Debug {
			log.Println(u.err)
		}
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

func (u *Urllib) Debug() *Urllib {
	u.config.Debug = true
	return u
}

func (u *Urllib) DisguisedIP(ip string) *Urllib {
	u.SetHeaderMap(map[string]string{
		"X-Forwarded-For":  ip,
		"X-Forwarded-Host": ip,
		"X-Client-IP":      ip,
		"X-remote-IP":      ip,
		"X-remote-addr":    ip,
		"True-Client-IP":   ip,
		"Client-IP":        ip,
		"X-Real-IP":        ip,
	})
	return u
}

func (u *Urllib) RandDisguisedIP() *Urllib {
	ip := lib.RandomIp()
	u.SetHeaderMap(map[string]string{
		"X-Forwarded-For":  ip,
		"X-Forwarded-Host": ip,
		"X-Client-IP":      ip,
		"X-remote-IP":      ip,
		"X-remote-addr":    ip,
		"True-Client-IP":   ip,
		"Client-IP":        ip,
		"X-Real-IP":        ip,
	})
	return u
}

func (u *Urllib) basicRules() {
	// Url Params
	distUrl, err := lib.BuildURLParams(u.url, u.querys)
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

	if u.config.Debug {
		fmt.Println("Url: ", u.req.URL)
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

func (u *Urllib) KeepAlives() *Urllib {
	u.disableKeepAlives = false
	return u
}

func (u *Urllib) body() (*http.Response, error) {
	u.basicRules()

	if u.config.Proxy == nil {
		u.config.Proxy = http.ProxyFromEnvironment
	}

	if u.config.Transport == nil {

		u.config.Transport = &http.Transport{
			TLSClientConfig:     u.config.TLSClientConfig,
			TLSHandshakeTimeout: u.config.ConnectTimeout,
			DisableKeepAlives:   u.disableKeepAlives,
			MaxIdleConnsPerHost: 100,
			Proxy:               u.config.Proxy,
			//Dial:                lib.SetTimeoutDialer(u.config.ConnectTimeout, u.config.ReadWriteTimeout),
		}
	} else {
		if t, ok := u.config.Transport.(*http.Transport); ok {
			if t.TLSClientConfig == nil {
				t.TLSClientConfig = u.config.TLSClientConfig
			}
			if t.Proxy == nil {
				t.Proxy = u.config.Proxy
			}
			//if t.Dial == nil {
			//	t.Dial = lib.SetTimeoutDialer(u.config.ConnectTimeout, u.config.ReadWriteTimeout)
			//}
		}
	}

	initCookie()
	jar := globalCookie

	client := http.Client{
		Transport: u.config.Transport,
		Jar:       jar,
		Timeout:   u.config.ConnectTimeout,
	}

	if u.config.CheckRedirect != nil {
		client.CheckRedirect = u.config.CheckRedirect
	}

	if u.err != nil {
		if u.config.Debug {
			log.Println(u.err)
		}
		return nil, u.err
	}

	u.setBody()
	if u.config.Gzip {
		u.setGzip()
	}

	return client.Do(u.req)
}

func (u *Urllib) byteOriginal() (int, []byte, error) {
	body, err := u.body()
	if err != nil {
		if u.config.Debug {
			log.Println(err)
		}
		return 0, nil, err
	}
	defer body.Body.Close()

	all, err := ioutil.ReadAll(body.Body)
	if err != nil {
		if u.config.Debug {
			log.Println(err)
		}
		return 0, nil, err
	}

	return body.StatusCode, all, nil
}

func (u *Urllib) byte() (int, []byte, error) {
	body, err := u.body()
	if err != nil {
		if u.config.Debug {
			log.Println(err)
		}
		return 0, nil, err
	}
	defer body.Body.Close()
	all, err := ioutil.ReadAll(body.Body)
	if err != nil {
		if u.config.Debug {
			log.Println(err)
		}
		return 0, all, err
	}
	if len(all) == 0 {
		return body.StatusCode, all, nil
	}
	// GZIP
	if body.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(bytes.NewReader(all))
		if err != nil {
			if u.config.Debug {
				log.Println(err)
			}
			return 0, nil, err
		}
		all, err = ioutil.ReadAll(reader)
		if err != nil {
			if u.config.Debug {
				log.Println(err)
			}
			return 0, nil, err
		}
	}

	// 旋转木马
	contentType := body.Header.Get("Content-Type")
	all, err = lib.AutomaticTranscoding(contentType, all)
	if err != nil {
		if u.config.Debug {
			log.Println(err)
		}
	}
	return body.StatusCode, all, err
}

func (u *Urllib) Body() (*http.Response, error) {
	return u.body()
}

func (u *Urllib) Byte() (int, []byte, error) {
	return u.byte()
}

func (u *Urllib) ByteOriginal() (int, []byte, error) {
	return u.byteOriginal()
}

// 拥有重尝 版本
func (u *Urllib) BodyRetry(retry int) (body *http.Response, err error) {
	if retry == 0 {
		retry = 3
	}

	for i := 0; i < retry; i++ {
		body, err = u.body()
		if err != nil {
			if i == 3 {
				switch {
				case i == 0:
					time.Sleep(time.Second * time.Duration(lib.Random(1, 5)))
				case i == 1:
					time.Sleep(time.Second * time.Duration(lib.Random(8, 10)))
				default:
					time.Sleep(time.Second * time.Duration(lib.Random(10, 20)))
				}
			} else {
				time.Sleep(time.Second * time.Duration(lib.Random(1, 5)))
			}
			continue
		}
		return body, err
	}
	return body, err
}

func (u *Urllib) ByteOriginalRetry(retry int, code int) (statusCode int, body []byte, err error) {
	if retry == 0 {
		retry = 3
	}
	for i := 0; i < retry; i++ {
		statusCode, body, err = u.byteOriginal()
		if code == 0 {
			if err != nil {
				if i == 3 {
					switch {
					case i == 0:
						time.Sleep(time.Second * time.Duration(lib.Random(1, 5)))
					case i == 1:
						time.Sleep(time.Second * time.Duration(lib.Random(8, 10)))
					default:
						time.Sleep(time.Second * time.Duration(lib.Random(10, 20)))
					}
				} else {
					time.Sleep(time.Second * time.Duration(lib.Random(1, 5)))
				}
				continue
			}
		} else {
			if err != nil || statusCode != code {
				if i == 3 {
					switch {
					case i == 0:
						time.Sleep(time.Second * time.Duration(lib.Random(1, 5)))
					case i == 1:
						time.Sleep(time.Second * time.Duration(lib.Random(8, 10)))
					default:
						time.Sleep(time.Second * time.Duration(lib.Random(10, 20)))
					}
				} else {
					time.Sleep(time.Second * time.Duration(lib.Random(1, 5)))
				}
				continue
			}
		}

		return statusCode, body, err
	}
	return statusCode, body, err
}

func (u *Urllib) ByteRetry(retry int, code int) (statusCode int, body []byte, err error) {
	if retry == 0 {
		retry = 3
	}
	for i := 0; i < retry; i++ {
		statusCode, body, err = u.byte()
		if code == 0 {
			if err != nil {
				if i == 3 {
					switch {
					case i == 0:
						time.Sleep(time.Second * time.Duration(lib.Random(1, 5)))
					case i == 1:
						time.Sleep(time.Second * time.Duration(lib.Random(8, 10)))
					default:
						time.Sleep(time.Second * time.Duration(lib.Random(10, 20)))
					}
				} else {
					time.Sleep(time.Second * time.Duration(lib.Random(1, 5)))
				}
				continue
			}
		} else {
			if err != nil || statusCode != code {
				if i == 3 {
					switch {
					case i == 0:
						time.Sleep(time.Second * time.Duration(lib.Random(1, 5)))
					case i == 1:
						time.Sleep(time.Second * time.Duration(lib.Random(8, 10)))
					default:
						time.Sleep(time.Second * time.Duration(lib.Random(10, 20)))
					}
				} else {
					time.Sleep(time.Second * time.Duration(lib.Random(1, 5)))
				}
				continue
			}
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

func (u *Urllib) FromJson(r interface{}) error {
	statusCode, body, err := u.byte()
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return errors.New(string(body))
	}

	return json.Unmarshal(body, r)
}

func (u *Urllib) FromJsonByCode(r interface{}, code int) error {
	statusCode, body, err := u.byte()
	if err != nil {
		return err
	}

	if statusCode != code {
		return errors.New(string(body))
	}

	return json.Unmarshal(body, r)
}
