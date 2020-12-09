package urllib

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/dollarkillerx/fasthttp"
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

	method string
	params *fasthttp.Args
	querys url.Values
	header map[string]string

	data   []byte
	files  map[string]string
	config Config

	err error
}

type Config struct {
	UserAgent       string
	Timeout         time.Duration
	TLSClientConfig *tls.Config
	Redirect        int // 最大跳转次数
	Gzip            bool
}

var defaultConfig = Config{
	UserAgent:       "Urllib Gold FastHttp",
	Timeout:         30 * time.Second,
	Redirect:        3,
	Gzip:            false,
	TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
}

func getBase(tagUrl, method string) *Urllib {
	base := &Urllib{
		params: &fasthttp.Args{},
		querys: map[string][]string{},
		files:  map[string]string{},
		header: map[string]string{},
		url:    tagUrl,
		err:    nil,
		config: defaultConfig,
		method: method,
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
	u.header[k] = v
	return u
}

func (u *Urllib) SetHeaderMap(headMap map[string]string) *Urllib {
	for k, v := range headMap {
		u.header[k] = v
	}
	return u
}

func (u *Urllib) SetUserAgent(ua string) *Urllib {
	u.config.UserAgent = ua
	return u
}

func (u *Urllib) RandUserAgent() *Urllib {
	u.config.UserAgent = lib.ReptileGetUserAgent()
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
	u.config.Timeout = timeout
	return u
}

func (u *Urllib) SetTLSClientConfig(config *tls.Config) *Urllib {
	u.config.TLSClientConfig = config
	return u
}

//func (u *Urllib) SetHTTPVersion(version string) *Urllib {
//	major, minor, ok := http.ParseHTTPVersion(version)
//	if ok {
//		u.req.Proto = version
//		u.req.ProtoMajor = major
//		u.req.ProtoMinor = minor
//	}
//
//	return u
//}
//
//func (u *Urllib) SetCookie(cookie *http.Cookie) *Urllib {
//	u.req.Header.Add("Cookie", cookie.String())
//	return u
//}
//
//func (u *Urllib) SetTransport(transport http.RoundTripper) *Urllib {
//	u.config.Transport = transport
//	return u
//}
//
//func (u *Urllib) SetCheckRedirect(redirect func(req *http.Request, via []*http.Request) error) *Urllib {
//	u.config.CheckRedirect = redirect
//	return u
//}
//
//func (u *Urllib) SetAuth(user, password string) *Urllib {
//	u.req.SetBasicAuth(user, password)
//	return u
//}
//
//func (u *Urllib) SetProxy(proxy func(*http.Request) (*url.URL, error)) *Urllib {
//	u.config.Proxy = proxy
//	return u
//}
//
//func (u *Urllib) PostFile(fromKey, filePath string) *Urllib {
//	u.files[fromKey] = filePath
//	return u
//}

//func (u *Urllib) setBody() {
//	if len(u.data) != 0 {
//		u.data = ioutil.NopCloser(bytes.NewBuffer(u.data))
//	}
//}

func (u *Urllib) SetBody(body []byte) *Urllib {
	u.data = body
	return u
}

func (u *Urllib) SetJson(body []byte) *Urllib {
	if body != nil {
		u.data = body
		u.header["Content-Type"] = "application/json"
	}
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
		u.header["Content-Type"] = "application/json"
	}
	return u
}

// if form
//func (u *Urllib) setFormParams(Forms url.Values) {
//	data := Forms.Encode()
//	u.req.Body = ioutil.NopCloser(strings.NewReader(data))
//	u.req.ContentLength = int64(len(data))
//}

func (u *Urllib) setGzip() *Urllib {
	if u.data == nil || len(u.data) == 0 {
		return u
	}

	data, err := lib.GZipData(u.data)
	if err != nil {
		u.err = err
		return u
	}

	u.SetHeader("Content-Encoding", "gzip")
	u.data = data
	return u
}

func (u *Urllib) NoRedirect() *Urllib {
	u.config.Redirect = 0
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
	u.url = parse.String()
}

func (u *Urllib) body() (result *http.Response, err error) {
	if u.err != nil {
		return result, u.err
	}

	u.basicRules()
	if u.config.Gzip {
		u.setGzip()
	}

	// FastHttp server
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // <- do not forget to release

	req.Header.SetMethod(u.method)
	req.SetRequestURI(u.url)
	if len(u.params.String()) > 0 {
		if _, err := u.params.WriteTo(req.BodyWriter()); err != nil {
			return result, u.err
		}
		u.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	}

	if len(u.data) > 0 {
		req.SetBody(u.data)
	}

	// set header
	for k, v := range u.header {
		req.Header.Set(k, v)
	}

	// ua
	req.Header.SetUserAgent(u.config.UserAgent)

	return u.doRequestFollowRedirectsBuffer(req, nil, u.url, &client, u.config.Redirect)
}

var client fasthttp.Client

func init() {
	client = fasthttp.Client{
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}

}

type clientDoer interface {
	Do(req *fasthttp.Request, resp *fasthttp.Response) error
	DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration) error
}

func (u *Urllib) clientGetURL(dst []byte, url string, c clientDoer) (response *http.Response, err error) {
	req := fasthttp.AcquireRequest()

	buffer, err := u.doRequestFollowRedirectsBuffer(req, dst, url, c, 3)

	fasthttp.ReleaseRequest(req)
	return buffer, err
}

func (u *Urllib) doRequestFollowRedirectsBuffer(req *fasthttp.Request, dst []byte, url string, c clientDoer, redirectsCount int) (result *http.Response, err error) {
	resp := fasthttp.AcquireResponse()
	bodyBuf := resp.BodyBuffer()
	resp.KeepBodyBuffer = true
	oldBody := bodyBuf.B
	bodyBuf.B = dst
	var body []byte
	result, err = u.doRequestFollowRedirects(req, resp, url, redirectsCount, c)
	body = bodyBuf.B
	bodyBuf.B = oldBody
	resp.KeepBodyBuffer = false
	fasthttp.ReleaseResponse(resp)
	body = body
	if err != nil {
		return nil, err
	}

	return result, err
}

func (u *Urllib) doRequestFollowRedirects(req *fasthttp.Request, resp *fasthttp.Response, url string, maxRedirectsCount int, c clientDoer) (result *http.Response, err error) {
	redirectsCount := 0
	var statusCode int

	for {
		req.SetRequestURI(url)
		if err := req.ParseURI(); err != nil {
			return nil, err
		}

		if err = c.DoTimeout(req, resp, u.config.Timeout); err != nil {
			break
		}
		statusCode = resp.Header.StatusCode()
		if !fasthttp.StatusCodeIsRedirect(statusCode) {
			break
		}

		redirectsCount++
		if redirectsCount > maxRedirectsCount {
			err = fasthttp.ErrTooManyRedirects
			break
		}
		location := resp.Header.Peek(fasthttp.HeaderLocation)
		if len(location) == 0 {
			err = fasthttp.ErrMissingLocation
			break
		}
		url = lib.RedirectUrl(url, string(location))
	}

	result = &http.Response{
		StatusCode: resp.StatusCode(),
		Header:     http.Header{},
	}

	resp.Header.VisitAll(func(key, value []byte) {
		result.Header.Add(string(key), string(value))
	})

	result.Body = ioutil.NopCloser(bytes.NewReader(resp.Body()))
	result.ContentLength = int64(len(resp.Body()))
	return result, nil
}

func (u *Urllib) byte() (int, []byte, error) {
	body, err := u.body()
	if err != nil {
		log.Println(err)
		return 0, nil, err
	}
	defer body.Body.Close()
	var all []byte
	// GZIP
	if body.Header.Get("Content-Encoding") == "gzip" {
		r, err := gzip.NewReader(body.Body)
		if err != nil {
			log.Println(err)
			return 0, nil, err
		}
		defer r.Close()

		reader := bufio.NewReader(r)
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				log.Println(err)
				break
			}
			all = append(all, line...)
		}

	} else {
		all, err = ioutil.ReadAll(body.Body)
		if err != nil {
			log.Println(err)
			return 0, all, err
		}
		if len(all) == 0 {
			return body.StatusCode, all, nil
		}
	}

	// 旋转木马
	contentType := body.Header.Get("Content-Type")
	all, err = lib.AutomaticTranscoding(contentType, all)
	if err != nil {
		log.Println(err)
	}
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
			time.Sleep(time.Second * time.Duration(lib.Random(1, 5)))
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
			time.Sleep(time.Second * time.Duration(lib.Random(1, 5)))
			continue
		}
		return statusCode, body, err
	}
	return statusCode, body, err
}
