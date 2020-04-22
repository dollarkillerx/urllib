package urllib

import (
	"net/http"
	"net/url"
	"time"
)

type urlType int

const (
	UA = "UrlLib By DollarKillerx"

	GET urlType = iota
	POST
	PUT
	DELETE
)

type urllib struct {
	url     string
	typ     urlType
	timeout time.Duration

	ua        string
	noTls     bool
	httpProxy string
	req       *http.Request
	resp      *http.Response

	params  url.Values
	header  map[string]string
	cookies []*http.Cookie
}

func Get(url string) *urllib {
	return &urllib{
		url:     url,
		typ:     GET,
		ua:      UA,
		timeout: time.Second * 3,
		params:  map[string][]string{},
		header:  map[string]string{},
		cookies: []*http.Cookie{},
	}
}

func Post(url string) *urllib {
	return &urllib{
		url:     url,
		typ:     GET,
		ua:      UA,
		timeout: time.Second * 3,
		params:  map[string][]string{},
		header:  map[string]string{},
		cookies: []*http.Cookie{},
	}
}

func Put(url string) *urllib {
	return &urllib{
		url:     url,
		typ:     GET,
		ua:      UA,
		timeout: time.Second * 3,
		params:  map[string][]string{},
		header:  map[string]string{},
		cookies: []*http.Cookie{},
	}
}

func Delete(url string) *urllib {
	return &urllib{
		url:     url,
		typ:     GET,
		ua:      UA,
		timeout: time.Second * 3,
		params:  map[string][]string{},
		header:  map[string]string{},
		cookies: []*http.Cookie{},
	}
}

func (u *urllib) SetTimeOut(timeout time.Duration) *urllib {
	u.timeout = timeout
	return u
}

func (u *urllib) SetUserAgent(ua string) *urllib {
	u.ua = ua
	return u
}

func (u *urllib) RandUserAgent() *urllib {
	u.ua = ReptileGetUserAgent()
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
func (u *urllib) NoTLS(b bool) *urllib {
	u.noTls = b
	return u
}

func (u *urllib) HttpProxy(url string) *urllib {
	u.httpProxy = url
	return u
}

func (u *urllib) Body() (http.Response, error) {
	client := http.Client{}

	return client.Do()
}
